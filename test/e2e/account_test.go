package e2e_test

import (
	"encoding/json"
	"github.com/ecommerce-api/internal/handler"
	"github.com/ecommerce-api/pkg/constant"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/security"
	"github.com/ecommerce-api/test/mocks/service_mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"testing"
	"time"
)

var accountService = service_mocks.AccountServiceMock{Mock: mock.Mock{}}

var accountHandler = handler.NewAccountHandler(&accountService)

func popData(role string) entity.User {
	password, _ := security.HashPassword("localhost")
	email := "test@mail.com"

	userData := entity.User{
		BaseEntity: entity.BaseEntity{ID: 1},
		Username:   "test",
		Email:      email,
		Password:   password,
		Role:       role,
	}

	security.PayloadData = &security.Payload{
		UserID:   userData.ID,
		Role:     userData.Role,
		Username: userData.Username,
		Email:    userData.Email,
		Exp:      time.Now().Add(360).Unix(),
	}

	return userData
}

func TestPing(t *testing.T) {
	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	req := httptest.NewRequest("GET", "/ping", nil)

	// Perform the request plain with the app,
	// the second argument is a request latency
	// (set to -1 for no latency)
	resp, _ := app.Test(req, 1)

	// Verify, if the status code is as expected
	assert.Equalf(t, 200, resp.StatusCode, "get HTTP status 200")
}

func TestAccountOwner(t *testing.T) {

	userOwner := popData(constant.RoleOwner)

	t.Run("Test SqlDB profile", func(t *testing.T) {
		accountService.Mock.On("Profile", userOwner.ID).Return(userOwner, nil)

		app.Get("/profile", accountHandler.Profile)

		req := httptest.NewRequest("GET", "/profile", nil)
		resp, _ := app.Test(req, 10)

		// Verify, if the status code is as expected
		assert.Equalf(t, 200, resp.StatusCode, "get HTTP status 200")

		//verify body
		assert.NotZero(t, resp.ContentLength)

		//parse body
		bodyData := make([]byte, resp.ContentLength)
		_, _ = resp.Body.Read(bodyData)
		var respBody map[string]entity.User
		_ = json.Unmarshal(bodyData, &respBody)

		//assert equal email
		assert.Equal(t, userOwner.Email, respBody["data"].Email)

	})

}
