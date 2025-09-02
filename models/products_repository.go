package models

import (
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(offset, limit int, category string, priceLt *float64) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.Model(&Product{}).Preload("Variants").Preload("Category")

	// Apply filters
	if category != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").Where("categories.name = ?", category)
	}

	if priceLt != nil {
		query = query.Where("products.price < ?", *priceLt)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated products
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
