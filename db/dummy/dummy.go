package dummy

import (
	"fmt"
	"github.com/ecommerce-api/pkg/constant"
	. "github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/helper"
	"github.com/ecommerce-api/pkg/security"
	"gorm.io/gorm"
	"strconv"
)

var password, _ = security.HashPassword("123456")

type Dummy interface {
	ClearData()
	CreateUsers(count int) []User
	CreateShops(count int, users []User) []Shop
	CreateProduct(count int, shops []Shop) []Product
	CreateWarehouses(count int, shops []Shop) []Warehouse
	CreateWarehousesInventory(warehouses []Warehouse, products []Product) []WarehouseInventory
}

type dummy struct {
	db *gorm.DB
}

func (d dummy) ClearData() {
	fmt.Println("clear all data")
	d.db.Exec("DELETE FROM stock_locks")
	d.db.Exec("DELETE FROM order_warehouse_allocations")
	d.db.Exec("DELETE FROM product_transfer_warehouses")
	d.db.Exec("DELETE FROM warehouse_inventories")
	d.db.Exec("DELETE FROM warehouses")
	d.db.Exec("DELETE FROM shops")
	d.db.Exec("DELETE FROM order_details")
	d.db.Exec("DELETE FROM payments")
	d.db.Exec("DELETE FROM orders")
	d.db.Exec("DELETE FROM products")
	d.db.Exec("DELETE FROM profiles")
	d.db.Exec("DELETE FROM users")
}

func (d dummy) CreateUsers(count int) []User {
	fmt.Print("create users ---> STARTED")
	var users []User

	if count > 1 {
		for i := 0; i < count; i++ {
			var u = User{
				Username: "owner",
				Email:    "owner@gmail.com",
				Password: password,
				Role:     constant.RoleOwner,
				Profile: &Profile{
					Name:       "Owner",
					Phone:      "11111111",
					Address:    "Yogyakarta",
					City:       "Yogyakarta",
					State:      "Daerah Istimewa Yogyakarta",
					PostalCode: "1234",
					Country:    "Indonesia",
				},
			}

			users = append(users, u)
		}
	}

	var defaultUsers = []User{
		{
			Username: "owner",
			Email:    "owner@gmail.com",
			Password: password,
			Role:     constant.RoleOwner,
			Profile: &Profile{
				Name:       "Owner",
				Phone:      "11111111",
				Address:    "Yogyakarta",
				City:       "Yogyakarta",
				State:      "Daerah Istimewa Yogyakarta",
				PostalCode: "1234",
				Country:    "Indonesia",
			},
		},
		{
			Username: "customer",
			Email:    "customer@gmail.com",
			Password: password,
			Role:     constant.RoleCustomer,
			Profile: &Profile{
				Name:       "Customer",
				Phone:      "22222222",
				Address:    "Magelang",
				City:       "Magelang",
				State:      "Jawa Tengah",
				PostalCode: "4321",
				Country:    "Indonesia",
			},
		},
	}

	users = append(users, defaultUsers...)

	if err := d.db.Create(&users).Error; err != nil {
		panic(err)
	}

	fmt.Println("-------> DONE")

	return users
}

func (d dummy) CreateShops(count int, users []User) []Shop {
	fmt.Print("create shops ---> STARTED")
	if count < 1 {
		count = 1
	}

	var shops []Shop

	for _, user := range users {
		for i := 0; i < count; i++ {
			var Shop = Shop{
				Name:        fmt.Sprintf("Shop %s", user.Username),
				Description: fmt.Sprintf("Shop %s Description", user.Username),
				Address:     "Jakarta",
				Phone:       "1234567" + strconv.Itoa(i),
				UserID:      user.ID, //need to set manually
			}

			shops = append(shops, Shop)
		}
	}

	if err := d.db.Model(&Shop{}).Create(shops).Error; err != nil {
		panic(err)
	}

	fmt.Println(" ------> DONE")

	return shops
}

func (d dummy) CreateProduct(count int, shops []Shop) []Product {
	fmt.Print("create product ---> STARTED")
	if count < 1 {
		count = 1
	}

	var products []Product

	for _, shop := range shops {
		for i := 0; i < count; i++ {
			var index = i + 1
			var p = Product{
				ShopID:      shop.ID, //need to set manually
				Name:        fmt.Sprintf("Product  %s - %d", shop.Name, index),
				Description: fmt.Sprintf("Desc Product %s - %d ", shop.Name, index),
				Price:       10000 * float64(index),
				SKU:         helper.GenRandomString(8),
				ImageURL:    fmt.Sprintf("https://picsum.photos/id/%d/200/300", index),
			}

			products = append(products, p)
		}
	}

	if err := d.db.Model(&Product{}).Create(&products).Error; err != nil {
		panic(err)
	}

	fmt.Println(" -----> DONE")

	return products
}

func (d dummy) CreateWarehouses(count int, shops []Shop) []Warehouse {
	fmt.Print("create warehouse ---> STARTED")

	if count < 1 {
		count = 1
	}

	var wareHouses []Warehouse

	for _, shop := range shops {
		for i := 1; i <= count; i++ {
			var index = i + 1
			wh := Warehouse{
				Name:     fmt.Sprintf("Warehouse %s - %d", shop.Name, index),
				IsActive: true,
				Location: "Yogyakarta",
				ShopID:   shop.ID,
				UserID:   shop.UserID,
			}

			wareHouses = append(wareHouses, wh)
		}
	}

	if err := d.db.Model(&Warehouse{}).Create(&wareHouses).Error; err != nil {
		panic(err)
	}

	fmt.Println(" ----> DONE")
	return wareHouses
}

func (d dummy) CreateWarehousesInventory(warehouses []Warehouse, products []Product) []WarehouseInventory {
	fmt.Print("create warehouse inventory ---> STARTED")

	var warehouseInventories []WarehouseInventory

	for _, warehouse := range warehouses {
		for i, product := range products {
			wi := WarehouseInventory{
				ProductID:   product.ID,
				WarehouseID: warehouse.ID,
				Quantity:    10 * (i + 1),
			}
			warehouseInventories = append(warehouseInventories, wi)
		}
	}
	if err := d.db.Model(&WarehouseInventory{}).Create(&warehouseInventories).Error; err != nil {
		panic(err)
	}

	fmt.Println(" ----> DONE")

	return warehouseInventories

}

func NewDummy(db *gorm.DB) Dummy {
	return dummy{db: db}
}
