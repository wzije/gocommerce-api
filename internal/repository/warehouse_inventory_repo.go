package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/ecommerce-api/pkg/constant"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/repository"
	"github.com/ecommerce-api/pkg/security"
	"gorm.io/gorm"
)

type warehouseRepository struct {
	db *gorm.DB
}

func (w *warehouseRepository) CreateProductInventory(ctx context.Context, warehouseID uint64, productID uint64, quantity int) error {

	var warehouseInventory entity.WarehouseInventory

	if err := w.db.WithContext(ctx).
		Model(&entity.WarehouseInventory{}).
		Where("product_id=?", productID).
		First(&warehouseInventory).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			var inventory = entity.WarehouseInventory{
				WarehouseID: warehouseID,
				ProductID:   productID,
				Quantity:    quantity,
			}

			if err := w.db.Model(&entity.WarehouseInventory{}).
				Create(&inventory).
				Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil

}

func (w *warehouseRepository) GetWarehouseByUser(ctx context.Context, userID uint64) (*[]entity.Warehouse, error) {
	var warehouses []entity.Warehouse
	err := w.db.WithContext(ctx).Model(&entity.Warehouse{}).
		Preload("WarehouseInventory").
		Find(&warehouses, "user_id=?", userID).Error
	if err != nil {
		return nil, err
	}
	return &warehouses, nil
}

func (w *warehouseRepository) GetWarehouseByShop(ctx context.Context, shopID uint64) (*[]entity.Warehouse, error) {
	var warehouses []entity.Warehouse
	err := w.db.WithContext(ctx).Model(&entity.Warehouse{}).
		Preload("WarehouseInventory").
		First(&warehouses, "shop_id=?", shopID).Error
	if err != nil {
		return nil, err
	}
	return &warehouses, nil
}

func (w *warehouseRepository) SelectWarehouse(ctx context.Context, productID uint64, shopID uint64, requiredQuantity int) (*entity.Warehouse, error) {

	var warehouse entity.Warehouse

	err := w.db.Model(&entity.Warehouse{}).
		Joins("JOIN warehouse_inventories ON warehouse_inventories.warehouse_id = warehouses.id").
		Where("warehouses.shop_id=? AND warehouse_inventories.product_id = ? "+
			"AND warehouse_inventories.quantity >= ? AND warehouses.is_active IS TRUE AND warehouses.user_id = ?",
			shopID, productID, requiredQuantity, security.PayloadData.UserID).
		First(&warehouse).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no warehouse found") //Tidak ada gudang dengan stok yang mencukupi
		}
		return nil, err
	}

	return &warehouse, nil
}

func (w *warehouseRepository) GetWarehouseByID(ctx context.Context, warehouseID uint64) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	err := w.db.WithContext(ctx).Model(&warehouse).
		First(&warehouse, "id=? AND user_id=?", warehouseID, security.PayloadData.UserID).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (w *warehouseRepository) UpdateWarehouseStatus(ctx context.Context, warehouseID uint64, isActive bool) error {
	warehouse, err := w.GetWarehouseByID(ctx, warehouseID)
	if err != nil {
		return err
	}
	warehouse.IsActive = isActive
	return w.db.WithContext(ctx).Model(&warehouse).Save(warehouse).Error
}

func (w *warehouseRepository) GetAvailableStock(ctx context.Context, productID uint64, shopID uint64) (int, error) {

	var totalQuantity int

	err := w.db.WithContext(ctx).Model(&entity.WarehouseInventory{}).
		Select("sum(warehouse_inventories.quantity) as total_quantity").
		Joins("join warehouses on warehouses.id = warehouse_inventories.warehouse_id").
		Where("warehouses.shop_id = ? AND warehouse_inventories.product_id = ? AND warehouses.user_id = ?",
			shopID, productID, security.PayloadData.UserID).
		Scan(&totalQuantity).Error

	if err != nil {
		return 0, err
	}

	return totalQuantity, nil

}

func (w *warehouseRepository) GetInventory(ctx context.Context, warehouseID, productID uint64) (*entity.WarehouseInventory, error) {
	var inventory entity.WarehouseInventory
	err := w.db.WithContext(ctx).Model(&inventory).
		Joins("LEFT JOIN warehouses ON warehouse_inventories.warehouse_id = warehouses.id").
		Where("warehouse_inventories.warehouse_id = ? AND warehouse_inventories.product_id = ? AND warehouses.user_id=?",
			warehouseID, productID, security.PayloadData.UserID).
		First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (w *warehouseRepository) ReduceStock(ctx context.Context, productID uint64, warehouseID uint64, quantity int) error {
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inventory entity.WarehouseInventory
		if err := tx.WithContext(ctx).Model(&entity.WarehouseInventory{}).
			Joins("LEFT JOIN warehouses ON (warehouse_inventories.warehouse_id = warehouses.id)").
			Where("warehouse_inventories.warehouse_id = ? AND warehouse_inventories.product_id = ?", warehouseID, productID).
			Where("warehouses.user_id = ?", security.PayloadData.UserID).
			First(&inventory).Error; err != nil {
			return err
		}

		if inventory.Quantity < quantity {
			return fmt.Errorf("insufficient stock to reduce")
		}

		inventory.Quantity -= quantity
		if err := tx.Save(&inventory).Error; err != nil {
			return err
		}

		return nil
	})
}

func (w *warehouseRepository) IncreaseStock(ctx context.Context, productID uint64, warehouseID uint64, quantity int) error {
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inventory entity.WarehouseInventory
		if err := tx.Model(&entity.WarehouseInventory{}).
			Joins("LEFT JOIN warehouses ON warehouse_inventories.warehouse_id = warehouses.id").
			Where("warehouse_inventories.warehouse_id = ? AND warehouse_inventories.product_id = ?", warehouseID, productID).
			Where("warehouses.user_id = ?", security.PayloadData.UserID).
			First(&inventory).Error; err != nil {
			return err
		} else {
			// Jika ada entri, tambahkan kuantitas yang ada
			inventory.Quantity += quantity
			if err := tx.Model(&entity.WarehouseInventory{}).
				Where("product_id=? AND warehouse_id=?", inventory.ProductID, warehouseID).
				Save(&inventory).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (w *warehouseRepository) TransferStock(ctx context.Context, sourceWarehouseID, destinationWarehouseID, productID uint64, quantity int) error {
	// Start a transaction
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check source warehouse inventory
		sourceInventory, err := w.GetInventory(ctx, sourceWarehouseID, productID)
		if err != nil {
			return err
		}

		if sourceInventory.Quantity < quantity {
			return errors.New("not enough stock in source warehouse")
		}

		// Deduct from source warehouse
		sourceInventory.Quantity -= quantity
		if err := tx.Model(&entity.WarehouseInventory{}).
			Where("id =?", sourceInventory.ID).
			Save(&sourceInventory).Error; err != nil {
			return err
		}

		// Add to destination warehouse
		destinationWarehouse, err := w.GetInventory(ctx, sourceWarehouseID, productID)
		if err != nil {
			return err
		}

		var destinationInventory entity.WarehouseInventory
		err = tx.
			Model(&entity.WarehouseInventory{}).
			Where("warehouse_id = ? AND product_id = ?", destinationWarehouse.ID, productID).
			First(&destinationInventory).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				destinationInventory = entity.WarehouseInventory{
					WarehouseID: destinationWarehouseID,
					ProductID:   productID,
					Quantity:    quantity,
				}
				if err := tx.Model(&entity.WarehouseInventory{}).
					Create(&destinationInventory).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			destinationInventory.Quantity += quantity
			if err := tx.Model(&entity.WarehouseInventory{}).
				Save(&destinationInventory).Error; err != nil {
				return err
			}
		}

		// Record the transfer
		productTransfer := entity.ProductTransferWarehouse{
			SourceWarehouseID:      sourceWarehouseID,
			DestinationWarehouseID: destinationWarehouseID,
			ProductID:              productID,
			Quantity:               quantity,
			Status:                 constant.ProductTransferComplete,
		}

		if err := tx.Model(&entity.ProductTransferWarehouse{}).
			Create(&productTransfer).Error; err != nil {
			return err
		}

		return nil
	})
}

func NewWarehouseInventoryRepository(db *gorm.DB) repository.WarehouseInventoryRepositoryInterface {
	return &warehouseRepository{
		db: db,
	}
}
