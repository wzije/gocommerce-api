package repository

import (
	"gorm.io/gorm"
)

type BaseRepositoryInterface[T any] interface {
	Fetch(offset int, limit int, q string, sort string, filter string) (*[]T, error)
	List() (*[]T, error)
	GetByID(id uint64) (*T, error)
	Store(item *T) (*T, error)
	Update(id uint64, updatedItem *T) (*T, error)
	Delete(id uint64) error
	TotalData() (int64, error)
	GetDB() *gorm.DB
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepositoryInterface[T] {
	return &baseRepository[T]{db: db}
}

func (repo *baseRepository[T]) GetDB() *gorm.DB {
	return repo.db
}

func (repo *baseRepository[T]) List() (*[]T, error) {
	var items []T
	var model T

	result := repo.db.Model(&model).Find(&items)

	return &items, result.Error
}

func (repo *baseRepository[T]) Fetch(offset int, limit int, q string, sort string, filter string) (*[]T, error) {
	var items []T

	query := repo.db

	if q != "" {
		query = repo.db.Where(
			repo.db.Where("name ILIKE ? ", "%"+q+"%"),
		)
	}

	if sort != "" {
		query = query.Order(sort)
	}

	if filter != "" {
		switch filter {
		case "IS_ACTIVE":
			query = query.Where("status = ?", true)
		case "IS_INACTIVE":
			query = query.Where("status = ?", false)
		}
	}

	if err := query.
		Limit(limit).
		Offset(offset).
		Find(&items).Error; err != nil {
		return nil, err
	}

	return &items, nil
}

func (repo *baseRepository[T]) Store(item *T) (*T, error) {
	result := repo.db.Model(&item).Create(&item)
	return item, result.Error
}

func (repo *baseRepository[T]) GetByID(id uint64) (*T, error) {
	var item T
	result := repo.db.Model(&item).First(&item, id)
	return &item, result.Error
}

func (repo *baseRepository[T]) Update(id uint64, updatedItem *T) (*T, error) {
	var item T
	result := repo.db.Model(&item).First(&item, id)
	if result.Error != nil {
		return nil, result.Error
	}
	result = repo.db.Model(&item).Model(&item).Updates(updatedItem)
	return &item, result.Error
}

func (repo *baseRepository[T]) Delete(id uint64) error {
	var item T
	result := repo.db.Model(&item).Delete(&item, id)
	return result.Error
}

func (repo *baseRepository[T]) TotalData() (int64, error) {
	var count int64
	var item T
	result := repo.db.Model(&item).Count(&count)

	return count, result.Error
}
