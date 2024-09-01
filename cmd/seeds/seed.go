package main

import (
	"github.com/ecommerce-api/db/dummy"
	"github.com/ecommerce-api/pkg/config"
	"log"
)

var db config.DB

func init() {
	log.Println("Execute Init...")
	config.Load()
	db = config.SqlDBLoad()
}

func main() {
	d := dummy.NewDummy(db.SqlDB())

	d.ClearData()

	users := d.CreateUsers(1)
	shops := d.CreateShops(1, users)
	products := d.CreateProduct(5, shops)
	warehouse := d.CreateWarehouses(1, shops)
	_ = d.CreateWarehousesInventory(warehouse, products)

}
