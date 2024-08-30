package security

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ecommerce-api/pkg/config"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/http"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type Payload struct {
	UserID   uint64
	Role     string
	Username string
	Email    string
	Exp      int64
}

var (
	JwtSigningMethod = jwt.SigningMethodHS256.Name
	PayloadData      = new(Payload)
)

func CreateToken(user *entity.User) (string, int64, error) {
	// Store token

	exp := time.Now().Add(time.Hour * 96).Unix()

	//set claim
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"iss":      "access",
		"aud":      "node service",
		"sub":      user.Username,
		"iat":      time.Now().Unix(),
		"exp":      exp,
		"nbf":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecretKey, err := ReadJWTSecretKey()

	if err != nil {
		return "", 0, err
	}

	result, err := token.SignedString(jwtSecretKey)

	if err != nil {

		return "", 0, errors.New("generate token failed")
	}

	return result, exp, nil
}

func validateSignedMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
	}

	jwtSecretKey, err := ReadJWTSecretKey()

	if err != nil {
		return nil, err
	}

	return jwtSecretKey, nil
}

func ParseToken(tokenString string) (*jwt.MapClaims, error) {
	claims := new(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, validateSignedMethod)
	if err != nil {
		return nil, err
	}

	var ok bool
	claims, ok = token.Claims.(*jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, http.ErrInvalidAuthToken
	}

	return claims, nil
}

func ParsePayload(ctx *fiber.Ctx) {
	token := ctx.Locals("jwtKey").(*jwt.Token)
	payloadPointer := token.Claims.(*jwt.MapClaims)
	payload := *payloadPointer

	PayloadData = &Payload{
		UserID:   uint64(payload["user_id"].(float64)),
		Username: getStringValue(payload, "username"),
		Email:    getStringValue(payload, "email"),
		Role:     getStringValue(payload, "role"),
		Exp:      int64(payload["exp"].(float64)),
	}

}

// ReadJWTSecretKey read secret key from file
func ReadJWTSecretKey() ([]byte, error) {
	file, err := os.Open(config.Config.CertPath)

	if err != nil {
		logrus.Fatal("Open file failed")
		return nil, err
	}

	defer func() {
		if err = file.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	jwtRead, err := ioutil.ReadAll(file)

	return jwtRead, nil
}

func getStringValue(data jwt.MapClaims, key string) string {
	var value string
	if data[key] != nil {
		value = data[key].(string)
	}
	return value
}
