package repository_test

import (
	"context"
	"github.com/ecommerce-api/internal/repository"
	. "github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/helper"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestProductRepository(t *testing.T) {
	db := SetupTestDB(t)
	productRepo := repository.NewProductRepository(db)
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		product := Product{
			Name:        "name " + strconv.Itoa(i),
			Description: "description" + strconv.Itoa(i),
			Price:       float64(100*i + 1),
			SKU:         helper.GenRandomString(10),
			ImageURL:    "",
			ShopID:      1,
		}

		err := db.WithContext(ctx).Model(&Product{}).Create(&product).Error
		assert.NoError(t, err)
	}

	//Test List
	products, err := productRepo.List(ctx)
	productList := *products
	assert.NoError(t, err)
	assert.Len(t, productList, 10)

	// Test GetByID
	retrievedProduct, err := productRepo.GetById(ctx, productList[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, productList[0].Name, retrievedProduct.Name)

}
