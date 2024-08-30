package repository_test

import (
	"github.com/ecommerce-api/pkg/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate all the models
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Profile{},
		&entity.Product{},
		&entity.Shop{},
		&entity.Warehouse{},
		&entity.Order{},
		&entity.OrderDetail{},
		&entity.OrderWarehouseAllocation{},
		&entity.WarehouseInventory{},
		&entity.StockLock{},
	)
	if err != nil {
		t.Fatalf("failed to migrate models: %v", err)
	}

	return db
}
