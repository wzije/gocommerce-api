package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ecommerce-api/internal/repository"
	"github.com/ecommerce-api/pkg/constant"
	"github.com/ecommerce-api/pkg/dto"
	"github.com/ecommerce-api/pkg/entity"
	repoPkg "github.com/ecommerce-api/pkg/repository"
	"github.com/ecommerce-api/pkg/security"
	"github.com/ecommerce-api/pkg/service"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

type orderService struct {
	task          *sync.WaitGroup
	productRepo   repoPkg.ProductRepositoryInterface
	orderRepo     repoPkg.OrderRepositoryInterface
	stockLockRepo repoPkg.StockLockRepositoryInterface
	warehouseRepo repoPkg.WarehouseInventoryRepositoryInterface
	paymentRepo   repoPkg.PaymentRepositoryInterface
}

func (s *orderService) MyListOrder(ctx context.Context) (*[]entity.Order, error) {
	return s.orderRepo.MyListOrder(ctx)
}

func NewOrderService(db *gorm.DB, task *sync.WaitGroup) service.OrderServiceInterface {

	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	stockLockRepo := repository.NewStockLockRepository(db)
	warehouseRepo := repository.NewWarehouseInventoryRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	return &orderService{
		task:          task,
		productRepo:   productRepo,
		orderRepo:     orderRepo,
		stockLockRepo: stockLockRepo,
		warehouseRepo: warehouseRepo,
		paymentRepo:   paymentRepo,
	}
}

// CheckoutOrder Checkout CreateOrder Example of service method to handle an order creation
func (s *orderService) CheckoutOrder(ctx context.Context, request *dto.OrderRequest) (*entity.Order, error) {

	for _, detail := range request.Details {

		totalStock, err := s.warehouseRepo.GetAvailableStock(ctx, detail.ProductID, request.ShopID)

		if err != nil {
			return nil, err
		}

		//check availability stock
		if totalStock < detail.Quantity {
			return nil, errors.New("stock is not available")
		}

	}

	//mapping details
	var details []entity.OrderDetail
	var totalAmount float64
	for _, item := range request.Details {

		product, err := s.productRepo.GetById(ctx, item.ProductID)

		if err != nil {
			return nil, err
		}

		totalPrice := product.Price * float64(item.Quantity)
		totalAmount += totalPrice

		details = append(details, entity.OrderDetail{
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			PricePerUnit: product.Price,
			//TotalPrice:   totalPrice,
		})
	}

	//mapping order
	now := time.Now()
	order := &entity.Order{
		UserID:          security.PayloadData.UserID, //customer
		ShopID:          request.ShopID,
		Date:            now,
		Amount:          totalAmount,
		ShippingDate:    now.AddDate(0, 0, 3),
		Status:          constant.TransactionPending,
		ShippingAddress: request.ShippingAddress,
		OrderDetail:     details,
	}

	// Create the order
	order, err := s.orderRepo.CreateOrder(ctx, order)

	if err != nil {
		return nil, err
	}

	s.task.Add(2)

	//LOCK STOCK
	go s.lockStock(ctx, order)

	//RELEASE AFTER 1 MINUTE NO ACTION
	dateTicker := time.NewTicker(20 * time.Second)
	go s.releaseStock(ctx, order.ID, true, dateTicker)

	return order, nil
}

func (s *orderService) PaymentOrder(ctx context.Context, request *dto.PaymentRequest) error {
	//create payment
	//release lock
	//set order to paid / awaiting shipment

	order, err := s.orderRepo.GetOrderById(ctx, request.OrderID)

	if err != nil {
		return err
	}

	amount := float64(0)
	for _, detail := range order.OrderDetail {
		amount += *detail.TotalPrice
	}

	if amount != request.Amount {
		return errors.New("amount is invalid")
	}

	payment := &entity.Payment{
		OrderID:       request.OrderID,
		PaymentMethod: request.PaymentMethod,
		Amount:        request.Amount,           //populate from order
		Status:        constant.PaymentComplete, //bypass to paid
	}

	err = s.paymentRepo.CreatePayment(ctx, payment)

	if err != nil {
		return err
	}

	s.task.Add(1)

	//release stock without go back
	dateTicker := time.NewTicker(1 * time.Second)
	go s.releaseStock(ctx, order.ID, false, dateTicker)

	//set transaction status to complete
	err = s.orderRepo.UpdateOrderStatus(ctx, order.ID, constant.TransactionComplete)

	if err != nil {
		return err
	}

	return nil
}

// lockStock selecting the warehouse
// simple condition: the warehouse selected if the required stock available
func (s *orderService) lockStock(ctx context.Context, order *entity.Order) {
	s.task.Done()
	fmt.Println("lock stock triggered order ")
	for _, item := range order.OrderDetail {

		wh, err := s.warehouseRepo.SelectWarehouse(ctx, item.ProductID, order.ShopID, item.Quantity)

		if err != nil {
			logrus.Error(err)
			return
		}

		lock := &entity.StockLock{
			OrderID:     order.ID,
			ProductID:   item.ProductID,
			WarehouseID: wh.ID,
			Quantity:    item.Quantity,
		}

		if err = s.stockLockRepo.LockStock(ctx, lock); err != nil {
			logrus.Error(err)
			return
		}

		if err = s.warehouseRepo.ReduceStock(context.Background(),
			lock.ProductID, lock.WarehouseID, lock.Quantity); err != nil {
			logrus.Error(err)
			return
		}

		info := fmt.Sprintf(
			"lock stock order: %d, product: %d, total: %d from warehouse %d",
			order.ID, lock.ProductID, lock.Quantity, lock.WarehouseID)
		fmt.Println(info)
		logrus.Info(info)

	}
}

func (s *orderService) releaseStock(ctx context.Context, orderID uint64, isBack bool, ticker *time.Ticker) {
	s.task.Done()

	<-ticker.C
	locks, err := s.stockLockRepo.GetStockLockByOrder(ctx, orderID)

	if err != nil {
		logrus.Error(err)
		return
	}

	if len(*locks) == 0 {
		fmt.Print("no lock stock available")
	}

	for _, lock := range *locks {
		info := fmt.Sprintf(
			"release stock order %d, product: %d, total: %d  from warehose %d ",
			orderID, lock.ProductID, lock.Quantity, lock.WarehouseID)
		fmt.Println(info)
		logrus.Info(info)

		if err := s.stockLockRepo.ReleaseStock(ctx, lock.ID); err != nil {
			logrus.Error(err)
			return
		}

		if isBack {
			if err := s.warehouseRepo.IncreaseStock(context.Background(), lock.ProductID, lock.WarehouseID, lock.Quantity); err != nil {
				logrus.Error(err)
				return
			}

			if err := s.orderRepo.UpdateOrderStatus(context.Background(), orderID, constant.TransactionCancelled); err != nil {
				logrus.Error(err)
				return
			}
		}

	}

}

func (s *orderService) ReleaseAllOldStock(ctx context.Context, t *time.Time) {
	//s.task.Done()

	locks, err := s.stockLockRepo.GetAllStockLockOlderThan(ctx, t)

	if err != nil {
		logrus.Error(err)
		fmt.Println("nothing to release")
		return
	}

	for _, lock := range *locks {
		info := fmt.Sprintf(
			"release stock - product: %d, total: %d  from warehose %d ",
			lock.ProductID, lock.Quantity, lock.WarehouseID)
		fmt.Println(info)
		logrus.Info(info)

		if err := s.stockLockRepo.ReleaseStock(ctx, lock.ID); err != nil {
			logrus.Error(err)
			return
		}

		if err := s.warehouseRepo.IncreaseStock(context.Background(), lock.ProductID, lock.WarehouseID, lock.Quantity); err != nil {
			logrus.Error(err)
			return
		}
	}

}
