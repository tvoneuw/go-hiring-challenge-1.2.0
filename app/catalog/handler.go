package catalog

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Total    int64     `json:"total"`
	Products []Product `json:"products"`
}

type Product struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type ProductsFetcher interface {
	GetAllProducts(offset, limit int, category string, priceLt *float64) ([]models.Product, int64, error)
}

type CatalogHandler struct {
	repo ProductsFetcher
}

func NewCatalogHandler(r ProductsFetcher) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	offset, limit := parsePaginationParams(r)

	// Parse filters
	category, priceLt := parseFilters(r)

	res, total, err := h.repo.GetAllProducts(offset, limit, category, priceLt)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map response
	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	// Return the products as a JSON response
	api.OKResponse(w, Response{
		Total:    total,
		Products: products,
	})
}
