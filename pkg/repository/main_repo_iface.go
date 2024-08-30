package repository

import (
	"context"
	"github.com/ecommerce-api/pkg/dto"
	"github.com/ecommerce-api/pkg/entity"
	"time"
)

// UserRepositoryInterface include profile
type UserRepositoryInterface interface {
	Register(user entity.User) (*entity.User, error)
	ByEmail(email string) (*entity.User, error)
	ByEmailOrPhone(email string, phone string) (*entity.User, error)
	ByPhone(phone string) (*entity.User, error)
	Fetch(offset int, limit int, q string, sort string, filter string) (*[]entity.User, error)
	List() (*[]entity.User, error)
	ById(id uint64) (*entity.User, error)
	Store(user *entity.User) (*entity.User, error)
	Update(id uint64, user dto.UserRequest) (*entity.User, error)
	Delete(id uint64) error
	TotalData() (int64, error)
}

type ProductRepositoryInterface interface {
	List(ctx context.Context) (*[]entity.Product, error)
	ListWithStock(ctx context.Context) (*[]dto.ProductResponse, error)
	GetById(ctx context.Context, id uint64) (*entity.Product, error)
	GetByIdWithStock(ctx context.Context, id uint64) (*dto.ProductResponse, error)
}

type ShopRepositoryInterface interface {
	BaseRepositoryInterface[entity.Shop]
}

type WarehouseRepositoryInterface interface {
}

type OrderWarehouseAllocationRepositoryInterface interface {
}

type OrderRepositoryInterface interface {
	MyListOrder(ctx context.Context) (*[]entity.Order, error)
	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	GetOrderById(ctx context.Context, id uint64) (*entity.Order, error)
	UpdateOrder(ctx context.Context, order *entity.Order) error
	DeleteOrder(ctx context.Context, id uint64) error
	UpdateOrderStatus(ctx context.Context, id uint64, status string) error
}

type OrderDetailRepositoryInterface interface {
	CreateOrderDetail(ctx context.Context, detail *entity.OrderDetail) error
	GetOrderDetailsByOrderId(ctx context.Context, orderId uint64) ([]entity.OrderDetail, error)
	DeleteOrderDetail(ctx context.Context, id uint64) error
}

type WarehouseInventoryRepositoryInterface interface {
	GetWarehouseByUser(ctx context.Context, userID uint64) (*[]entity.Warehouse, error)
	GetWarehouseByShop(ctx context.Context, shopID uint64) (*[]entity.Warehouse, error)
	GetWarehouseByID(ctx context.Context, warehouseID uint64) (*entity.Warehouse, error)
	UpdateWarehouseStatus(ctx context.Context, warehouseID uint64, isActive bool) error
	GetInventory(ctx context.Context, warehouseID, productID uint64) (*entity.WarehouseInventory, error)
	TransferStock(ctx context.Context, sourceWarehouseID, destinationWarehouseID, productID uint64, quantity int) error
	GetAvailableStock(ctx context.Context, productId uint64, shopId uint64) (int, error)
	SelectWarehouse(ctx context.Context, productId uint64, shopId uint64, requiredQuantity int) (*entity.Warehouse, error)
	ReduceStock(ctx context.Context, productId uint64, warehouseId uint64, quantity int) error
	IncreaseStock(ctx context.Context, productId uint64, warehouseId uint64, quantity int) error
}

type StockLockRepositoryInterface interface {
	GetStockLockByOrder(ctx context.Context, orderId uint64) (*[]entity.StockLock, error)
	GetStockLockByOrderAndProduct(ctx context.Context, orderId uint64, productId uint64) (*entity.StockLock, error)
	LockStock(ctx context.Context, lock *entity.StockLock) error
	ReleaseStock(ctx context.Context, lockId uint64) error
	GetTotalLockedStock(ctx context.Context, orderId uint64, productId uint64) (int, error)
	GetAllStockLockOlderThan(ctx context.Context, t *time.Time) (*[]entity.StockLock, error)
}

type PaymentRepositoryInterface interface {
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	GetPaymentByOrderId(ctx context.Context, orderId uint64) (*entity.Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentId uint64, status string) error
}
