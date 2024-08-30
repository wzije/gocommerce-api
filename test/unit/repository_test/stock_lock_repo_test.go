package repository_test

import (
	"context"
	"github.com/ecommerce-api/internal/repository"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStockLockRepository(t *testing.T) {
	db := SetupTestDB(t)
	stockLockRepo := repository.NewStockLockRepository(db)

	t.Run("GetTotalLockedStock", func(t *testing.T) {
		stockLock := &entity.StockLock{
			OrderID:     1,
			ProductID:   1,
			WarehouseID: 1,
			Quantity:    10,
		}

		_ = stockLockRepo.LockStock(context.Background(), stockLock)

		lockedStock, err := stockLockRepo.GetTotalLockedStock(context.Background(), stockLock.OrderID, stockLock.ProductID)
		assert.NoError(t, err)
		assert.Equal(t, 10, lockedStock)
	})

	t.Run("lockStock", func(t *testing.T) {
		stockLock := &entity.StockLock{
			OrderID:     1,
			ProductID:   1,
			WarehouseID: 1,
			Quantity:    10,
		}

		err := stockLockRepo.LockStock(context.Background(), stockLock)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, stockLock.ID)
	})

	t.Run("releaseStock", func(t *testing.T) {
		stockLock := &entity.StockLock{
			OrderID:     1,
			ProductID:   1,
			WarehouseID: 1,
			Quantity:    10,
		}
		_ = stockLockRepo.LockStock(context.Background(), stockLock)

		err := stockLockRepo.ReleaseStock(context.Background(), stockLock.ID)
		assert.NoError(t, err)

		var releasedStockLock entity.StockLock
		db.First(&releasedStockLock, stockLock.ID)
		assert.NotNil(t, releasedStockLock.ReleasedAt)
	})

}
