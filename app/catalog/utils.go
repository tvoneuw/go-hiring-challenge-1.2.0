package catalog

import (
	"net/http"
	"strconv"
)

func parsePaginationParams(r *http.Request) (offset int, limit int) {
	// Default values
	offset = 0
	limit = 10

	// Parse offset
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			switch {
			case v <= 0:
				offset = 0
			default:
				offset = v
			}
		}
	}

	// Parse limit
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			switch {
			case v < 1:
				limit = 1
			case v > 100:
				limit = 100
			default:
				limit = v
			}
		}
	}

	return offset, limit
}

func parseFilters(r *http.Request) (category string, priceLt *float64) {
	category = r.URL.Query().Get("category")

	if p := r.URL.Query().Get("price_lt"); p != "" {
		if v, err := strconv.ParseFloat(p, 64); err == nil {
			priceLt = &v
		}
	}

	return category, priceLt
}
