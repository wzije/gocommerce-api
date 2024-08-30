package config

import (
	"fmt"
	"github.com/ecommerce-api/pkg/helper"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// DB public interface
type DB interface {
	SqlDB() *gorm.DB
}

// db struct accessing by internal
type db struct {
	slqDB *gorm.DB
}

// SqlDB Pgsql GetPgsql Get implement struct
func (d db) SqlDB() *gorm.DB {
	return d.slqDB
}

// SqlDBLoad logic for connection
func SqlDBLoad() DB {

	dbUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		helper.GetEnv("DB_HOST", "localhost"),
		helper.GetEnv("DB_USERNAME", "root"),
		helper.GetEnv("DB_PASSWORD", "localhost"),
		helper.GetEnv("DB_DATABASE", "database"),
		helper.GetEnv("DB_PORT", "3306"))

	fmt.Println("SqlDBLoad URL:", dbUrl)

	log.Println("connecting database...")

	database, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		logrus.Panic("database connection failed")
	}

	return &db{
		slqDB: database,
	}
}

func CustomDB(dialect gorm.Dialector) DB {
	dbase, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return &db{
		slqDB: dbase,
	}
}
