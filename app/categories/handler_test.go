package categories_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

type mockCategoryRepo struct {
	categories []models.Category
	err        error
}

func (m *mockCategoryRepo) GetAllCategories() ([]models.Category, error) {
	return m.categories, m.genericErr
}

func TestHandleGetCategories(t *testing.T) {
	categoriesSet := []models.Category{
		{Code: "clothing", Name: "Clothing"},
		{Code: "shoes", Name: "Shoes"},
	}
	tests := []struct {
		name       string
		repo       *mockCategoryRepo
		wantStatus int
		wantCount  int
	}{
		{
			name: "success",
			repo: &mockCategoryRepo{
				categories: categoriesSet,
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name: "repository_failure",
			repo: &mockCategoryRepo{
				genericErr: &mockError{},
			},
			wantStatus: http.StatusInternalServerError,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := categories.NewCategoriesHandler(tt.repo)

			req := httptest.NewRequest("GET", "/categories", nil)
			rec := httptest.NewRecorder()

			handler.HandleGet(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode, "unexpected status code")

			if tt.wantStatus == http.StatusOK {
				var resp categories.Response
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				assert.Equal(t, tt.wantCount, len(resp.Categories))
			}
		})
	}
}
