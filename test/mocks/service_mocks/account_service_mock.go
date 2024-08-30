package service_mocks

import (
	"github.com/ecommerce-api/pkg/dto"
	. "github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/http"
	"github.com/ecommerce-api/pkg/security"
	"github.com/ecommerce-api/pkg/service"
	"github.com/stretchr/testify/mock"
)

type AccountServiceMock struct {
	Mock mock.Mock
}

func (a *AccountServiceMock) Register(request *dto.AuthRegisterRequest) (*User, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountServiceMock) Login(request *dto.AuthAccessTokenRequest) (*dto.AuthAccessTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountServiceMock) RefreshToken(query *http.RequestQuery) error {
	//TODO implement me
	panic("implement me")
}

func (a *AccountServiceMock) Profile() (*User, error) {
	userId := security.PayloadData.UserID

	args := a.Mock.Called(userId)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	result := args.Get(0).(User)

	return &result, nil
}

func (a *AccountServiceMock) UpdateProfile(request dto.UserProfileRequest) (*User, error) {
	//TODO implement me
	panic("implement me")
}

func NewAccountServiceMock() service.AccountServiceInterface {
	return &AccountServiceMock{Mock: mock.Mock{}}
}
