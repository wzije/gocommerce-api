package service_mocks

import (
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/http"
	"github.com/stretchr/testify/mock"
)

type StaffServiceMock struct {
	Mock mock.Mock
}

func (s *StaffServiceMock) Fetch(query *http.RequestQuery) (*[]entity.User, *http.Pagination, error) {
	args := s.Mock.Called(query)

	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}

	if args.Get(1) == nil {
		return nil, nil, args.Error(2)
	}

	result := args.Get(0).([]entity.User)

	pagination := args.Get(1).(http.Pagination)

	return &result, &pagination, nil
}

func (s *StaffServiceMock) List() (*[]entity.User, error) {
	args := s.Mock.Called()

	if args.Get(0) != nil {
		return nil, args.Error(1)
	}

	result := args.Get(0).([]entity.User)

	return &result, nil
}

func (s *StaffServiceMock) ById(id uint64) (*entity.User, error) {
	args := s.Mock.Called(id)
	if args.Get(0) != nil {
		return nil, args.Error(1)
	}

	result := args.Get(0).(entity.User)

	return &result, nil

}

func (s *StaffServiceMock) Delete(id uint64) error {
	//TODO implement me
	panic("implement me")
}

func (s *StaffServiceMock) ChangeStatus(id uint64, status string) error {
	//TODO implement me
	panic("implement me")
}
