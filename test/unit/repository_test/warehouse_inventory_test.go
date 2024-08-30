package repository_test

import (
	"context"
	"github.com/ecommerce-api/internal/repository"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWarehouseInventoryRepository(t *testing.T) {
	db := SetupTestDB(t)
	warehouseRepo := repository.NewWarehouseInventoryRepository(db)

	t.Run("GetAvailableStock", func(t *testing.T) {

		warehouseInventory := &entity.WarehouseInventory{
			ProductID:   1,
			WarehouseID: 1,
			Quantity:    100,
		}

		db.Create(warehouseInventory)

		stock, err := warehouseRepo.GetAvailableStock(context.Background(), 1, 1)
		assert.NoError(t, err)
		assert.True(t, stock == 100)

		isNotAvailable, err := warehouseRepo.GetAvailableStock(context.Background(), 1, 1)
		assert.NoError(t, err)
		assert.False(t, isNotAvailable == 100)
	})

	t.Run("ReduceStock", func(t *testing.T) {
		warehouseInventory := &entity.WarehouseInventory{
			ProductID:   1,
			WarehouseID: 1,
			Quantity:    100,
		}
		db.Create(warehouseInventory)

		err := warehouseRepo.ReduceStock(context.Background(), 1, 1, 10)
		assert.NoError(t, err)

		var updatedInventory entity.WarehouseInventory
		db.First(&updatedInventory, warehouseInventory.ID)
		assert.Equal(t, 90, updatedInventory.Quantity)
	})

	t.Run("IncreaseStock", func(t *testing.T) {
		warehouseInventory := &entity.WarehouseInventory{
			ProductID:   1,
			WarehouseID: 1,
			Quantity:    100,
		}
		db.Create(warehouseInventory)

		err := warehouseRepo.IncreaseStock(context.Background(), 1, 1, 10)
		assert.NoError(t, err)

		var updatedInventory entity.WarehouseInventory
		db.First(&updatedInventory, warehouseInventory.ID)
		assert.Equal(t, 110, updatedInventory.Quantity)
	})
}
