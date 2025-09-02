package categories

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoriesFetcher interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(c *models.Category) error
	GetCategoryByCode(code string) (*models.Category, error)
}

type CategoriesHandler struct {
	repo CategoriesFetcher
}

func NewCategoriesHandler(r CategoriesFetcher) *CategoriesHandler {
	return &CategoriesHandler{
		repo: r,
	}
}

func (h *CategoriesHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	res, err := h.repo.GetAllCategories()

	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map Response
	categories := make([]Category, len(res))
	for i, c := range res {
		categories[i] = Category{
			Code: c.Code,
			Name: c.Name,
		}
	}

	// Return the products as a JSON response
	api.OKResponse(w, Response{
		Categories: categories,
	})
}

func (h *CategoriesHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var input models.Category
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate category input
	if err := validateCategoryInput(&input); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check code uniqueness
	category, err := h.repo.GetCategoryByCode(input.Code)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if category != nil {
		api.ErrorResponse(w, http.StatusConflict, "category with this code already exists")
		return
	}

	// Create effective category
	if err := h.repo.CreateCategory(&input); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Response
	api.OKResponse(w, Category{
		Code: input.Code,
		Name: input.Name,
	})
}
