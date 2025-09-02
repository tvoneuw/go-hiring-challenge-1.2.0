package details

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category string    `json:"category"`
	Variants []Variant `json:"variants"`
}

type Variant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type ProductFetcher interface {
	GetProductByCode(code string) (*models.Product, error)
}

type ProductDetailsHandler struct {
	repo ProductFetcher
}

func NewProductDetailsHandler(r ProductFetcher) *ProductDetailsHandler {
	return &ProductDetailsHandler{
		repo: r,
	}
}

func (h *ProductDetailsHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/catalog/")

	res, err := h.repo.GetProductByCode(code)
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	// Map response
	variants := make([]Variant, len(res.Variants))
	// Apply variant price inheritance
	for i, v := range res.Variants {
		if v.Price.IsZero() {
			v.Price = res.Price
		}

		variants[i] = Variant{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: v.Price.InexactFloat64(),
		}
	}

	response := Response{
		Code:     res.Code,
		Price:    res.Price.InexactFloat64(),
		Category: res.Category.Name,
		Variants: variants,
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
