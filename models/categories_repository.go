package models

import (
	"errors"

	"gorm.io/gorm"
)

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

func (r *CategoriesRepository) GetAllCategories() ([]Category, error) {
	var categories []Category

	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoriesRepository) CreateCategory(c *Category) error {
	return r.db.Create(c).Error
}

func (r *CategoriesRepository) GetCategoryByCode(code string) (*Category, error) {
	var category Category

	if err := r.db.Where("categories.code = ?", code).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}
