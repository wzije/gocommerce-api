package repository_test

import (
	"context"
	"github.com/ecommerce-api/internal/repository"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderRepository(t *testing.T) {
	db := SetupTestDB(t)
	orderRepo := repository.NewOrderRepository(db)

	t.Run("CreateOrder", func(t *testing.T) {
		order := &entity.Order{
			UserID:          1,
			Status:          "Pending",
			Amount:          100.00,
			ShippingAddress: "123 Street",
		}

		_, err := orderRepo.CreateOrder(context.Background(), order)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, order.ID)
	})

	t.Run("GetOrderById", func(t *testing.T) {
		order := &entity.Order{
			UserID:          1,
			Status:          "Pending",
			Amount:          100.00,
			ShippingAddress: "123 Street",
		}
		_, _ = orderRepo.CreateOrder(context.Background(), order)

		fetchedOrder, err := orderRepo.GetOrderById(context.Background(), order.ID)
		assert.NoError(t, err)
		assert.Equal(t, order.ID, fetchedOrder.ID)
		assert.Equal(t, "Pending", fetchedOrder.Status)
	})

	t.Run("UpdateOrder", func(t *testing.T) {
		order := &entity.Order{
			UserID:          1,
			Status:          "Pending",
			Amount:          100.00,
			ShippingAddress: "123 Street",
		}
		_, _ = orderRepo.CreateOrder(context.Background(), order)

		order.Status = "Completed"
		err := orderRepo.UpdateOrder(context.Background(), order)
		assert.NoError(t, err)

		updatedOrder, _ := orderRepo.GetOrderById(context.Background(), order.ID)
		assert.Equal(t, "Completed", updatedOrder.Status)
	})

	t.Run("DeleteOrder", func(t *testing.T) {
		order := &entity.Order{
			UserID:          1,
			Status:          "Pending",
			Amount:          100.00,
			ShippingAddress: "123 Street",
		}

		_, _ = orderRepo.CreateOrder(context.Background(), order)

		err := orderRepo.DeleteOrder(context.Background(), order.ID)
		assert.NoError(t, err)

		deletedOrder, err := orderRepo.GetOrderById(context.Background(), order.ID)
		assert.Error(t, err)
		assert.Nil(t, deletedOrder)
	})
}
