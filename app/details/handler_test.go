package details_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/details"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

// Fake repository for testing
type mockProductsRepository struct{}

func (m *mockProductsRepository) GetProductByCode(code string) (*models.Product, error) {
	if code != "PROD001" {
		return nil, fmt.Errorf("not found")
	}
	return &models.Product{
		Code:  "PROD001",
		Price: decimal.NewFromFloat(10.99),
		Category: models.Category{
			Name: "Clothing",
		},
		Variants: []models.Variant{
			{Name: "Variant A", SKU: "SKU001A", Price: decimal.NewFromFloat(11.99)},
			{Name: "Variant B", SKU: "SKU001B", Price: decimal.Zero},
		},
	}, nil
}

func TestHandleGetProductDetails(t *testing.T) {
	handler := details.NewProductDetailsHandler(&mockProductsRepository{})

	tests := []struct {
		name             string
		code             string
		price            float64
		category         string
		expectedStatus   int
		expectedVariants []struct {
			Name  string
			Price float64
		}
	}{
		{
			name:           "TC 1 : existing product",
			code:           "PROD001",
			price:          10.99,
			category:       "Clothing",
			expectedStatus: http.StatusOK,
			expectedVariants: []struct {
				Name  string
				Price float64
			}{
				{"Variant A", 11.99},
				{"Variant B", 10.99}, // inherited price
			},
		},
		{
			name:           "TC 2 : nonexistent product",
			code:           "PROD999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/catalog/"+tt.code, nil)
			w := httptest.NewRecorder()

			handler.HandleGet(w, req)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Code == http.StatusOK {
				var resp details.Response

				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				// Check Products
				if resp.Code != tt.code {
					t.Errorf("expected code %s, got %s", tt.code, resp.Code)
				}

				if resp.Price != tt.price {
					t.Errorf("expected price %f, got %f", tt.price, resp.Price)
				}

				if resp.Category != tt.category {
					t.Errorf("expected category %s, got %s", tt.category, resp.Category)
				}

				// Check Variants
				if len(resp.Variants) != len(tt.expectedVariants) {
					t.Fatalf("expected %d variants, got %d", len(tt.expectedVariants), len(resp.Variants))
				}

				for i, v := range resp.Variants {
					if v.Name != tt.expectedVariants[i].Name {
						t.Errorf("variant name mismatch: expected %s, got %s", tt.expectedVariants[i].Name, v.Name)
					}
					if v.Price != tt.expectedVariants[i].Price {
						t.Errorf("variant price mismatch: expected %.2f, got %.2f", tt.expectedVariants[i].Price, v.Price)
					}
				}
			}
		})
	}
}
