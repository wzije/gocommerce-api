package main

import (
	"fmt"
	"github.com/ecommerce-api/internal/repository"
	"github.com/ecommerce-api/pkg/config"
	"github.com/ecommerce-api/pkg/constant"
	. "github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/helper"
	"github.com/ecommerce-api/pkg/security"
	"github.com/sirupsen/logrus"
	"log"
)

var db config.DB

func clearAllData() {
	db.SqlDB().Exec("DELETE FROM stock_locks")
	db.SqlDB().Exec("DELETE FROM order_warehouse_allocations")
	db.SqlDB().Exec("DELETE FROM product_transfer_warehouses")
	db.SqlDB().Exec("DELETE FROM warehouse_inventories")
	db.SqlDB().Exec("DELETE FROM warehouses")
	db.SqlDB().Exec("DELETE FROM shops")
	db.SqlDB().Exec("DELETE FROM order_details")
	db.SqlDB().Exec("DELETE FROM payments")
	db.SqlDB().Exec("DELETE FROM orders")
	db.SqlDB().Exec("DELETE FROM products")
	db.SqlDB().Exec("DELETE FROM profiles")
	db.SqlDB().Exec("DELETE FROM users")
}

func createUser() []User {
	fmt.Println("insert user data ---- START")
	password, err := security.HashPassword("123456")

	if err != nil {
		panic("failed to hash password")
	}

	var users = []User{
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

	var userList []User
	var repo = repository.NewUserRepository(db)
	for _, v := range users {
		if _, err := repo.Store(&v); err != nil {
			logrus.Fatal(err)
		}

		userList = append(userList, v)
	}

	fmt.Println("insert user data ---- DONE")
	return userList

}

func createShop(userId uint64) *Shop {
	fmt.Println("insert shop data ---- START")
	shop := Shop{
		Name:        "SHOP 1",
		Description: "shop 1 desc",
		Address:     "yogyakarta",
		Phone:       "1234567",
		UserID:      userId,
	}

	db.SqlDB().Model(&shop).Create(&shop)

	fmt.Println("insert shop data ---- DONE")
	return &shop
}

func createProducts(shopId uint64) []Product {
	fmt.Println("insert product data ---- START")
	var products = []Product{
		{
			ShopID:      shopId,
			Name:        "Product 1",
			Description: "Product 1 description",
			Price:       10000,
			SKU:         helper.GenRandomString(8),
			ImageURL:    "https://picsum.photos/id/9/200/300",
		},
		{
			ShopID:      shopId,
			Name:        "Product 2",
			Description: "Product 2 description",
			Price:       15000,
			SKU:         helper.GenRandomString(8),
			ImageURL:    "https://picsum.photos/id/160/200/300",
		},
		{
			ShopID:      shopId,
			Name:        "Product 3",
			Description: "Product 3 description",
			Price:       20000,
			SKU:         helper.GenRandomString(8),
			ImageURL:    "https://picsum.photos/id/96/200/300",
		},
		{
			ShopID:      shopId,
			Name:        "Product 4",
			Description: "Product 4 description",
			Price:       20000,
			SKU:         helper.GenRandomString(8),
			ImageURL:    "https://picsum.photos/id/30/200/300",
		},
		{
			ShopID:      shopId,
			Name:        "Product 5",
			Description: "Product 5 description",
			Price:       25000,
			SKU:         helper.GenRandomString(8),
			ImageURL:    "https://picsum.photos/id/21/200/300",
		},
	}

	var productsList []Product

	for _, v := range products {
		if err := db.SqlDB().Model(&Product{}).Create(&v).Error; err != nil {
			fmt.Println(err.Error())
		}
		productsList = append(productsList, v)

	}
	fmt.Println("insert product data ----- DONE")

	return productsList

}

func createWareHouse(shopId uint64, userId uint64, products []Product) *[]Warehouse {
	fmt.Println("insert warehouse data ---- START")

	var wareHouses []Warehouse

	for i := 1; i <= 2; i++ {
		wh := Warehouse{
			Name:     fmt.Sprintf("Warehouse %d", i),
			IsActive: true,
			Location: "Yogyakarta",
			ShopID:   shopId,
			UserID:   userId,
		}

		db.SqlDB().Model(&Warehouse{}).Create(&wh)

		wareHouses = append(wareHouses, wh)
	}

	fmt.Println("insert warehouse inventory data ---- START")
	var warehouseInventory []WarehouseInventory
	for index, product := range products {

		whIndex := 0
		if index%2 != 0 {
			whIndex = 1
		}

		wi := WarehouseInventory{
			ProductID:   product.ID,
			WarehouseID: wareHouses[whIndex].ID,
			Quantity:    10 * (index + 1),
		}

		db.SqlDB().Model(&WarehouseInventory{}).Create(&wi)

		warehouseInventory = append(warehouseInventory, wi)
	}

	fmt.Println("insert warehouse  data ---- DONE")
	return &wareHouses
}

func init() {
	log.Println("Execute Init...")
	config.Load()
	db = config.SqlDBLoad()
}

func main() {
	clearAllData()
	users := createUser()
	shop := createShop(users[0].ID)
	products := createProducts(shop.ID)
	createWareHouse(shop.ID, users[0].ID, products)
}
