package categories

import (
	"errors"

	"github.com/mytheresa/go-hiring-challenge/models"
)

var (
	ErrMissingCode = errors.New("code is required")
	ErrMissingName = errors.New("name is required")
)

// validateCategoryInput validates the fields of a category before creation.
func validateCategoryInput(c *models.Category) error {
	if c.Code == "" {
		return ErrMissingCode
	}
	if c.Name == "" {
		return ErrMissingName
	}
	return nil
}
