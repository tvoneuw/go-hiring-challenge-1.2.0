package catalog

import (
	"net/http"
	"strconv"
)

func parseOffsetLimit(r *http.Request) (offset int, limit int) {
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
