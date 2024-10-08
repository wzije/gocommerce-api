package service

import (
	"context"
	"github.com/ecommerce-api/pkg/dto"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/http"
	"time"
)

type UserServiceInterface interface {
	Fetch(query *http.RequestQuery) (*[]entity.User, *http.Pagination, error)
	List() (*[]entity.User, error)
	ById(id uint64) (*entity.User, error)
	ByEmail(email string) (*entity.User, error)
	Store(user entity.User) (*entity.User, error)
	Update(id uint64, user dto.UserRequest) (*entity.User, error)
	Delete(id uint64) error
}

type AccountServiceInterface interface {
	Register(request *dto.AuthRegisterRequest) (*entity.User, error)
	Login(request *dto.AuthAccessTokenRequest) (*dto.AuthAccessTokenResponse, error)
	RefreshToken(query *http.RequestQuery) error
	Profile() (*entity.User, error)
	UpdateProfile(request dto.UserProfileRequest) (*entity.User, error)
}

type ProductServiceInterface interface {
	List(ctx context.Context) (*[]entity.Product, error)
	ListWithStock(ctx context.Context) (*[]dto.ProductResponse, error)
	GetByID(ctx context.Context, id uint64) (*entity.Product, error)
	GetByIDWithStock(ctx context.Context, id uint64) (*dto.ProductResponse, error)
	GetAvailabilityStock(ctx context.Context, productID uint64, shopID uint64) (int, error)
}

type ShopServiceInterface interface {
	List(ctx context.Context) (*[]entity.Shop, error)
	GetByID(ctx context.Context, shopID uint64) (*entity.Shop, error)
	Orders(ctx context.Context, shopID uint64) (*[]entity.Order, error)
	Products(ctx context.Context, shopID uint64) (*[]entity.Product, error)
}

type OrderServiceInterface interface {
	MyOrders(ctx context.Context) (*[]entity.Order, error)
	MyCustomerOrders(ctx context.Context) (*[]entity.Order, error)
	CheckoutOrder(ctx context.Context, order *dto.OrderRequest) (*entity.Order, error)
	PaymentOrder(ctx context.Context, request *dto.PaymentRequest) error
}

type WarehouseServiceInterface interface {
	MyWarehouseList(ctx context.Context) (*[]entity.Warehouse, error)
	MyWarehouseByID(ctx context.Context, id uint64) (*entity.Warehouse, error)
	CreateWarehouse(ctx context.Context, warehouse *dto.WarehouseRequest) (*entity.Warehouse, error)
	CreateProductInventory(ctx context.Context, productID uint64, warehouseID uint64, quantity int) error
	IncreaseStock(ctx context.Context, productID uint64, warehouseID uint64, quantity int) error
	ReduceStock(ctx context.Context, productID uint64, warehouseID uint64, quantity int) error
	TransferStock(ctx context.Context, sourceWarehouseID, destinationWarehouseID, productID uint64, quantity int) error
	UpdateWarehouseStatus(ctx context.Context, warehouseID uint64, isActive bool) error
	ReleaseAllOldStock(ctx context.Context, t *time.Time)
}
