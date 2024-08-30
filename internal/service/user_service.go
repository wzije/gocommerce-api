package service

import (
	"fmt"
	"github.com/ecommerce-api/internal/repository"
	"github.com/ecommerce-api/pkg/config"
	"github.com/ecommerce-api/pkg/dto"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/http"
	. "github.com/ecommerce-api/pkg/repository"
	"github.com/ecommerce-api/pkg/service"
	"os"
	"strconv"
	"strings"
)

type userService struct {
	userRepository UserRepositoryInterface
}

func (o userService) Fetch(query *http.RequestQuery) (*[]entity.User, *http.Pagination, error) {

	var offset int

	if query.Page == 1 || query.Page == 0 {
		offset = 0
	} else {
		offset = (query.Page - 1) * query.PerPage
	}

	if query.PerPage == 0 {
		query.PerPage, _ = strconv.Atoi(os.Getenv("DATA_LIMIT"))
	}

	if query.Sort != "" {
		sorts := strings.Split(query.Sort, ".")
		query.Sort = fmt.Sprintf("%s %s", sorts[0], sorts[1])
	}

	result, err := o.userRepository.Fetch(offset, query.PerPage, query.Q, query.Sort, strings.ToUpper(query.Filter))

	var total int64 = 0
	total, _ = o.userRepository.TotalData()

	var count int64 = 0
	if result != nil {
		count = int64(len(*result))
	}

	pagination := http.Paginate(count, total, query)

	if result == nil {
		return nil, nil, err
	} else {
		return result, &pagination, nil
	}
}

func (o userService) List() (*[]entity.User, error) {

	result, err := o.userRepository.List()

	if result == nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (o userService) ById(id uint64) (*entity.User, error) {
	result, err := o.userRepository.ById(id)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (o userService) ByEmail(email string) (*entity.User, error) {
	result, err := o.userRepository.ByEmail(email)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (o userService) Store(outlet entity.User) (*entity.User, error) {
	result, err := o.userRepository.Store(&outlet)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (o userService) Update(id uint64, user dto.UserRequest) (*entity.User, error) {
	result, err := o.userRepository.Update(id, user)

	if err != nil {
		return nil, err
	}

	return result, nil

}

func (o userService) Delete(id uint64) error {
	err := o.userRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func NewUserService(db config.DB) service.UserServiceInterface {
	userRepository := repository.NewUserRepository(db)
	return &userService{userRepository: userRepository}
}
