package middleware

import (
	"github.com/ecommerce-api/pkg/http"
	"github.com/ecommerce-api/pkg/security"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

// AuthAPIMiddleware middleware
func AuthAPIMiddleware(ctx *fiber.Ctx) error {
	jwtSecretKey, err := security.ReadJWTSecretKey()
	if err != nil {
		return http.JsonError(ctx, err)
	}
	return jwtware.New(jwtware.Config{
		ContextKey:    "jwtKey",
		SigningKey:    jwtSecretKey,
		SigningMethod: security.JwtSigningMethod,
		TokenLookup:   "header:Authorization",
		Claims:        new(jwt.MapClaims),
		AuthScheme:    "Bearer",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			logrus.Error(http.ErrUnauthorized.Error())
			return http.Json(ctx, &http.Response{
				Code:    fiber.StatusUnauthorized,
				Message: http.ErrUnauthorized.Error(),
			})
		},
	})(ctx)
}
