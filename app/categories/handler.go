package categories

import (
	"encoding/json"
	"net/http"

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Categories: categories,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *CategoriesHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var input models.Category
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate category input
	if err := validateCategoryInput(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check code uniqueness
	category, err := h.repo.GetCategoryByCode(input.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if category != nil {
		http.Error(w, "category with this code already exists", http.StatusConflict)
		return
	}

	// Create effective category
	if err := h.repo.CreateCategory(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(Category{
		Code: input.Code,
		Name: input.Name,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
