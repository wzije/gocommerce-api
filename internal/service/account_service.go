package service

import (
	eror "errors"
	"github.com/ecommerce-api/internal/repository"
	"github.com/ecommerce-api/pkg/config"
	"github.com/ecommerce-api/pkg/constant"
	"github.com/ecommerce-api/pkg/dto"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/exception"
	"github.com/ecommerce-api/pkg/helper"
	"github.com/ecommerce-api/pkg/http"
	. "github.com/ecommerce-api/pkg/repository"
	"github.com/ecommerce-api/pkg/security"
	"github.com/ecommerce-api/pkg/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
)

type accountService struct {
	userRepository UserRepositoryInterface
}

func (a accountService) Profile() (*entity.User, error) {

	userID := security.PayloadData.UserID

	user, err := a.userRepository.ById(userID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a accountService) UpdateProfile(request dto.UserProfileRequest) (*entity.User, error) {

	user := dto.UserRequest{
		Name:     request.Name,
		Username: strings.ReplaceAll(helper.NormalizeString(request.Username), " ", "_"),
	}

	updated, err := a.userRepository.Update(security.PayloadData.UserID, user)

	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (a accountService) Register(request *dto.AuthRegisterRequest) (*entity.User, error) {
	//normalize email param
	request.Email = helper.NormalizeString(request.Email)

	//check the password should not be empty
	if strings.TrimSpace(request.Password) == "" {
		return nil, exception.ErrEmptyPassword
	}

	//get existing user by email
	if existingEmail, err := a.userRepository.ByEmail(request.Email); err != nil && !eror.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.TranslateErr(err)
	} else if existingEmail != nil {
		return nil, exception.New(fiber.StatusBadRequest, "the email is taken", nil)
	}

	//set user data
	username := strings.Split(request.Email, "@")
	password, err := security.HashPassword(request.Password)
	if err != nil {
		return nil, exception.TranslateErr(err)
	}

	//populate data
	var authRegister = entity.User{
		Username: username[0],
		Email:    request.Email,
		Password: password,
		Role:     constant.RoleCustomer,
	}

	//register it!!
	user, err := a.userRepository.Register(authRegister)

	//handle error when created
	if err != nil {
		return nil, exception.ErrUnprocessedEntity
	}

	//if no error send verification email / or no verification
	//if !utils.IsEnv("test") {
	//	if err = tasks.
	//		Queue.Add(tasks.
	//		NotifyVerificationRequest.
	//		WithArgs(context.Background(), authRegister),
	//	); err != nil {
	//		return nil, errors.New(fiber.StatusPreconditionFailed, "Send email job failed", err.Error())
	//	}
	//}

	//final success
	return user, nil
}

func (a accountService) Login(request *dto.AuthAccessTokenRequest) (*dto.AuthAccessTokenResponse, error) {

	//validate email param
	if !helper.IsEmail(request.Email) {
		return nil, exception.Message("param request does not email")
	}

	//get request by email
	user, err := a.userRepository.ByEmail(request.Email)

	//check existing request
	if err != nil {
		return nil, exception.Message("email not found")
	}

	//verify password
	if !security.VerifyPassword(request.Password, user.Password) {
		return nil, http.ErrInvalidCredentials
	}

	//create token
	token, exp, err := security.CreateToken(user)

	if err != nil {
		return nil, err
	}

	return &dto.AuthAccessTokenResponse{
		RefreshToken: "", //not implement yet
		TokenType:    "Bearer",
		ExpiresIn:    exp,
		AccessToken:  strings.Replace(token, "Bearer ", "", -1),
		Email:        user.Email,
		//UserId:       user.ID,
		//Username:     user.Username,
		//Role:         user.Role,
	}, nil
}

func (a accountService) RefreshToken(query *http.RequestQuery) error {
	return nil
}

func NewAccountService(db config.DB) service.AccountServiceInterface {
	userRepository := repository.NewUserRepository(db)
	return &accountService{
		userRepository: userRepository,
	}
}
