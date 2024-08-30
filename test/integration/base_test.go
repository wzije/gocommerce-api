package integration_test

import (
	"fmt"
	"github.com/ecommerce-api/internal/repository"
	"github.com/ecommerce-api/pkg/config"
	"github.com/ecommerce-api/pkg/entity"
	security2 "github.com/ecommerce-api/pkg/security"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RepositoryTestSuite struct {
	suite.Suite
	db config.DB
}

// SetupSuite Run before all tests start
func (s *RepositoryTestSuite) SetupSuite() {
	// Connect mock database
	fmt.Println("Setup suit")
	config.Load("test")
	s.db = config.SqlDBLoad()
}

// TearDownTest Run after each test finished
func (s *RepositoryTestSuite) TearDownTest() {
	// Truncate all tables
	fmt.Println("execute tear down test")
	s.db.SqlDB().Exec("DELETE FROM product_transfer_warehouses")
	s.db.SqlDB().Exec("DELETE FROM stock_locks")
	s.db.SqlDB().Exec("DELETE FROM order_warehouse_allocations")
	s.db.SqlDB().Exec("DELETE FROM warehouse_inventories")
	s.db.SqlDB().Exec("DELETE FROM warehouses")
	s.db.SqlDB().Exec("DELETE FROM shops")
	s.db.SqlDB().Exec("DELETE FROM payments")
	s.db.SqlDB().Exec("DELETE FROM order_details")
	s.db.SqlDB().Exec("DELETE FROM orders")
	s.db.SqlDB().Exec("DELETE FROM products")
	s.db.SqlDB().Exec("DELETE FROM profiles")
	s.db.SqlDB().Exec("DELETE FROM users")

}

// TearDownSuite Run after all tests has finished
func (s *RepositoryTestSuite) TearDownSuite() {
	/*	Run migrate down */
	fmt.Println("execute tear down suit")
	s.db.SqlDB().Exec("DELETE FROM product_transfer_warehouses")
	s.db.SqlDB().Exec("DELETE FROM stock_locks")
	s.db.SqlDB().Exec("DELETE FROM order_warehouse_allocations")
	s.db.SqlDB().Exec("DELETE FROM warehouse_inventories")
	s.db.SqlDB().Exec("DELETE FROM warehouses")
	s.db.SqlDB().Exec("DELETE FROM shops")
	s.db.SqlDB().Exec("DELETE FROM payments")
	s.db.SqlDB().Exec("DELETE FROM order_details")
	s.db.SqlDB().Exec("DELETE FROM orders")
	s.db.SqlDB().Exec("DELETE FROM products")
	s.db.SqlDB().Exec("DELETE FROM profiles")
	s.db.SqlDB().Exec("DELETE FROM users")
}

func (s *RepositoryTestSuite) registerUserData() *entity.User {
	authRepo := repository.NewUserRepository(s.db)

	//phone := "08120000001"
	email := "user@gmail.com"
	username := "user"
	password, _ := security2.HashPassword("localhost")

	user := entity.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	//create
	created, err := authRepo.Register(user)
	//created, err := repo.OrderStore(&user)
	if err != nil {
		fmt.Println(err)
	}

	return created

}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
