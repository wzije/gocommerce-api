package service

import (
	"context"
	"errors"
	"fmt"
	repositoryPkg "github.com/ecommerce-api/internal/repository"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/repository"
	"github.com/ecommerce-api/pkg/security"
	"github.com/ecommerce-api/pkg/service"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type warehouseService struct {
	inventoryRepo repository.WarehouseInventoryRepositoryInterface
	stockLockRepo repository.StockLockRepositoryInterface
}

func (w warehouseService) CreateProductInventory(ctx context.Context, productID uint64, warehouseID uint64, quantity int) error {
	return w.inventoryRepo.CreateProductInventory(ctx, productID, warehouseID, quantity)
}

func (w warehouseService) MyWarehouseList(ctx context.Context) (*[]entity.Warehouse, error) {
	return w.inventoryRepo.GetWarehouseByUser(ctx, security.PayloadData.UserID)
}

func (w warehouseService) IncreaseStock(ctx context.Context, productID uint64, warehouseID uint64, quantity int) error {
	return w.inventoryRepo.IncreaseStock(ctx, productID, warehouseID, quantity)
}

func (w warehouseService) ReduceStock(ctx context.Context, productID uint64, warehouseID uint64, quantity int) error {
	return w.inventoryRepo.ReduceStock(ctx, productID, warehouseID, quantity)
}

func (w warehouseService) TransferStock(ctx context.Context, sourceWarehouseID, destinationWarehouseID, productID uint64, quantity int) error {
	sourceWarehouse, err := w.inventoryRepo.GetWarehouseByID(ctx, sourceWarehouseID)
	if err != nil {
		return err
	}
	if !sourceWarehouse.IsActive {
		return errors.New("source warehouse is not active")
	}

	destinationWarehouse, err := w.inventoryRepo.GetWarehouseByID(ctx, destinationWarehouseID)
	if err != nil {
		return err
	}
	if !destinationWarehouse.IsActive {
		return errors.New("destination warehouse is not active")
	}

	return w.inventoryRepo.TransferStock(ctx, sourceWarehouseID, destinationWarehouseID, productID, quantity)
}

func (w warehouseService) ReleaseAllOldStock(ctx context.Context, t *time.Time) {
	//s.task.Done()

	locks, err := w.stockLockRepo.GetAllStockLockOlderThan(ctx, t)

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

		if err := w.stockLockRepo.ReleaseStock(ctx, lock.ID); err != nil {
			logrus.Error(err)
			return
		}

		if err := w.inventoryRepo.IncreaseStock(context.Background(), lock.ProductID, lock.WarehouseID, lock.Quantity); err != nil {
			logrus.Error(err)
			return
		}
	}
}

func (w warehouseService) UpdateWarehouseStatus(ctx context.Context, warehouseID uint64, isActive bool) error {
	return w.inventoryRepo.UpdateWarehouseStatus(ctx, warehouseID, isActive)
}

func NewWarehouseService(db *gorm.DB) service.WarehouseServiceInterface {
	inventoryRepo := repositoryPkg.NewWarehouseInventoryRepository(db)
	return &warehouseService{inventoryRepo: inventoryRepo}
}
