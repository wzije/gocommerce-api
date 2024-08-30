package repository

import (
	"context"
	"github.com/ecommerce-api/pkg/constant"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/repository"
	"gorm.io/gorm"
	"time"
)

type stockLockRepository struct {
	db *gorm.DB
}

func NewStockLockRepository(db *gorm.DB) repository.StockLockRepositoryInterface {
	return &stockLockRepository{
		db: db,
	}
}

func (r *stockLockRepository) LockStock(ctx context.Context, lock *entity.StockLock) error {
	lock.Status = constant.StockLock
	return r.db.WithContext(ctx).Create(lock).Error
}

func (r *stockLockRepository) ReleaseStock(ctx context.Context, lockId uint64) error {
	return r.db.WithContext(ctx).
		Model(&entity.StockLock{}).
		Where("id = ?", lockId).
		Update("released_at", time.Now()).
		Update("status", constant.StockRelease).Error
}

func (r *stockLockRepository) GetTotalLockedStock(ctx context.Context, orderId uint64, productId uint64) (int, error) {
	var totalLocked int

	err := r.db.WithContext(ctx).
		Model(&entity.StockLock{}).
		Where("order_id = ? AND product_id = ? AND released_at IS NULL", orderId, productId).
		Select("SUM(quantity)").Scan(&totalLocked).Error

	return totalLocked, err
}

func (r *stockLockRepository) GetStockLockByOrderAndProduct(ctx context.Context, orderId uint64, productId uint64) (*entity.StockLock, error) {
	lockStock := &entity.StockLock{}

	err := r.db.WithContext(ctx).
		Model(&entity.StockLock{}).
		First(lockStock, "order_id =? AND product_id=?", orderId, productId).Error

	return lockStock, err
}

func (r *stockLockRepository) GetStockLockByOrder(ctx context.Context, orderId uint64) (*[]entity.StockLock, error) {
	var lockStocks []entity.StockLock

	err := r.db.WithContext(ctx).
		Model(&entity.StockLock{}).
		Where("order_id = ? AND status = ?", orderId, constant.StockLock).
		Find(&lockStocks).Error

	return &lockStocks, err
}

func (r *stockLockRepository) GetAllStockLockOlderThan(ctx context.Context, t *time.Time) (*[]entity.StockLock, error) {
	var lockStocks []entity.StockLock

	err := r.db.WithContext(ctx).
		Model(&entity.StockLock{}).
		Where("locked_at < ? AND status = ?", t, constant.StockLock).
		Find(&lockStocks).Error

	return &lockStocks, err
}
