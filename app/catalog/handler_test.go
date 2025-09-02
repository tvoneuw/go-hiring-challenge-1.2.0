package catalog_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// Mock repository
type mockProductRepo struct {
	products []models.Product
	err      error
}

func (m *mockProductRepo) GetAllProducts(offset, limit int, category string, priceLt *float64) ([]models.Product, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}

	// Filter products
	var filtered []models.Product
	for _, p := range m.products {
		if category != "" && p.Category.Name != category {
			continue
		}
		if priceLt != nil && p.Price.InexactFloat64() >= *priceLt {
			continue
		}
		filtered = append(filtered, p)
	}

	total := int64(len(filtered))

	// Apply offset/limit
	start := offset
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

type mockError struct{}

func (e *mockError) Error() string {
	return "repo mock error"
}

func TestHandleGetCatalog(t *testing.T) {
	products := []models.Product{
		{
			Code:  "PROD001",
			Price: decimal.NewFromFloat(10.99),
			Category: models.Category{
				Name: "Clothing",
			},
		},
		{
			Code:  "PROD002",
			Price: decimal.NewFromFloat(5.50),
			Category: models.Category{
				Name: "Shoes",
			},
		},
		{
			Code:  "PROD003",
			Price: decimal.NewFromFloat(22.99),
			Category: models.Category{
				Name: "Accessories",
			},
		},
	}

	tests := []struct {
		name          string
		repo          *mockProductRepo
		query         string
		wantStatus    int
		wantCount     int
		wantTotal     int64
		expectedCodes []string
	}{
		{
			name: "repository_error",
			repo: &mockProductRepo{
				err: &mockError{},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "limit_1",
			repo: &mockProductRepo{
				products: products,
			},
			query:         "?limit=1",
			wantStatus:    http.StatusOK,
			wantCount:     1,
			wantTotal:     3,
			expectedCodes: []string{"PROD001"},
		},
		{
			name: "negative_limit",
			repo: &mockProductRepo{
				products: products,
			},
			query:         "?limit=-1",
			wantStatus:    http.StatusOK,
			wantCount:     1,
			wantTotal:     3,
			expectedCodes: []string{"PROD001"},
		},
		{
			name: "offset_1_limit_1",
			repo: &mockProductRepo{
				products: products,
			},
			query:         "?offset=1&limit=1",
			wantStatus:    http.StatusOK,
			wantCount:     1,
			wantTotal:     3,
			expectedCodes: []string{"PROD002"},
		},
		{
			name: "negative_offset",
			repo: &mockProductRepo{
				products: products,
			},
			query:         "?offset=-1",
			wantStatus:    http.StatusOK,
			wantCount:     3,
			wantTotal:     3,
			expectedCodes: []string{"PROD001", "PROD002", "PROD003"},
		},
		{
			name: "no_filter",
			repo: &mockProductRepo{
				products: products,
			},
			query:         "",
			wantStatus:    http.StatusOK,
			wantCount:     3,
			wantTotal:     3,
			expectedCodes: []string{"PROD001", "PROD002", "PROD003"},
		},
		{
			name: "category_filter",
			repo: &mockProductRepo{
				products: products,
			},
			query:         "?category=Shoes",
			wantStatus:    http.StatusOK,
			wantCount:     1,
			wantTotal:     1,
			expectedCodes: []string{"PROD002"},
		},
		{
			name: "priceLt_filter ",
			repo: &mockProductRepo{
				products: products,
			},
			query:         "?price_lt=11",
			wantStatus:    http.StatusOK,
			wantCount:     2,
			wantTotal:     2,
			expectedCodes: []string{"PROD001", "PROD002"},
		},
		{
			name: "all_filters",
			repo: &mockProductRepo{
				products: products,
			},
			query:         "?category=Shoes&price_lt=6",
			wantStatus:    http.StatusOK,
			wantCount:     1,
			wantTotal:     1,
			expectedCodes: []string{"PROD002"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := catalog.NewCatalogHandler(tt.repo)

			req := httptest.NewRequest("GET", "/catalog"+tt.query, nil)
			rec := httptest.NewRecorder()

			handler.HandleGet(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode, "Unexpected status code")

			var resp catalog.Response
			if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			assert.Equal(t, tt.wantCount, len(resp.Products), "Unexpected product count")

			for i, code := range tt.expectedCodes {
				assert.Equal(t, code, resp.Products[i].Code, "Unexpected product code at index %d", i)
			}

			assert.Equal(t, tt.wantTotal, resp.Total, "Unexpected total products")
		})
	}
}
