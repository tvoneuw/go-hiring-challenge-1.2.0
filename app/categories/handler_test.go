package categories_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

type mockCategoryRepo struct {
	categories  []models.Category
	existing    *models.Category
	genericErr  error
	specificErr error
}

func (m *mockCategoryRepo) GetAllCategories() ([]models.Category, error) {
	return m.categories, m.genericErr
}

func (m *mockCategoryRepo) CreateCategory(c *models.Category) error {
	return m.genericErr
}

func (m *mockCategoryRepo) GetCategoryByCode(code string) (*models.Category, error) {
	return m.existing, m.specificErr
}

type mockError struct{}

func (e *mockError) Error() string {
	return "repo mock error"
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

func TestHandlePostCategories(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repo       *mockCategoryRepo
		wantStatus int
	}{
		{
			name: "create_category_error",
			body: `{"code":"CAT001","name":"Bags"}`,
			repo: &mockCategoryRepo{
				genericErr: &mockError{},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "get_code_error",
			body: `{"code":"CAT001","name":"Bags"}`,
			repo: &mockCategoryRepo{
				specificErr: &mockError{},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "success",
			body:       `{"code":"CAT001","name":"Bags"}`,
			repo:       &mockCategoryRepo{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing_code",
			body:       `{"name":"Accessories"}`,
			repo:       &mockCategoryRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing_name",
			body:       `{"code":"acc"}`,
			repo:       &mockCategoryRepo{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate",
			body: `{"code":"CAT001","name":"Clothing"}`,
			repo: &mockCategoryRepo{
				existing: &models.Category{Code: "CAT001", Name: "Clothing"},
			},
			wantStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := categories.NewCategoriesHandler(tt.repo)

			req := httptest.NewRequest("POST", "/categories", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.HandlePost(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode, "Unexpected status code")
		})
	}
}
