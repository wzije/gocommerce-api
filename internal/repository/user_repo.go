package repository

import (
	"github.com/ecommerce-api/pkg/config"
	"github.com/ecommerce-api/pkg/dto"
	"github.com/ecommerce-api/pkg/entity"
	"github.com/ecommerce-api/pkg/exception"
	"github.com/ecommerce-api/pkg/helper"
	"github.com/ecommerce-api/pkg/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func (u userRepository) Fetch(offset int, limit int, q string, sort string, filter string) (*[]entity.User, error) {
	var users []entity.User

	query := u.db

	if q != "" {
		query = u.db.Where(
			u.db.Where("username ILIKE ? ", "%"+q+"%").
				Or("email ILIKE ? ", "%"+q+"%").
				Or("phone ILIKE ? ", "%"+q+"%").
				Or("name ILIKE ? ", "%"+q+"%"),
		)
	}

	if sort != "" {
		query = query.Order(sort)
	}

	switch filter {
	case "IS_ACTIVE":
		query = query.Where("status = ?", 1)
	case "IS_INACTIVE":
		query = query.Where("status = ?", 0)
	}

	if err := query.
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, exception.TranslateErr(err)
	}

	return &users, nil
}

func (u userRepository) List() (*[]entity.User, error) {
	var users []entity.User

	if err := u.db.
		Find(&users).
		Error; err != nil {
		return nil, exception.TranslateErr(err)
	}

	return &users, nil
}

func (u userRepository) ByEmail(email string) (*entity.User, error) {
	var user entity.User

	if result := u.db.
		Model(&entity.User{}).
		Where("email = ?", email).
		First(&user); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u userRepository) ByPhone(phone string) (*entity.User, error) {
	var user entity.User

	if result := u.db.
		Model(&entity.User{}).
		Where("phone = ?", phone).
		First(&user); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil

}

func (u userRepository) ByEmailOrPhone(email string, phone string) (*entity.User, error) {
	var user entity.User

	if result := u.db.
		Model(&entity.User{}).
		Where("email = ?", email).
		Or("phone = ?", phone).
		First(&user); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil

}

func (u userRepository) Delete(id uint64) error {
	if err := u.db.Where("id = ?", id).Delete(&entity.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (u userRepository) Update(id uint64, request dto.UserRequest) (*entity.User, error) {

	var user entity.User

	user.Username = helper.NormalizeString(request.Username)
	user.Email = request.Email
	user.Role = request.Role

	if err := u.db.Where("id = ?", id).
		First(&entity.User{}).
		Updates(&user).Error; err != nil {
		return nil, err
	}

	existingUser, err := u.ById(id)

	if err != nil {
		return &entity.User{}, err
	}

	return existingUser, nil
}

func (u userRepository) Store(user *entity.User) (*entity.User, error) {

	result := u.db.Create(&user)
	return user, result.Error
}

func (u userRepository) ById(id uint64) (*entity.User, error) {
	var user entity.User

	result := u.db.
		Preload("Profile").
		Where("id = ?", id).
		Find(&user)

	if user.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &user, result.Error
}

func (u userRepository) TotalData() (int64, error) {

	var count int64
	if err := u.db.
		Model(&entity.User{}).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (u userRepository) Register(user entity.User) (*entity.User, error) {

	if err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Create(&user).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &user, nil
}

func NewUserRepository(conn config.DB) repository.UserRepositoryInterface {
	return &userRepository{db: conn.SqlDB()}
}
