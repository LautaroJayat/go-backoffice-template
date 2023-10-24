package product_crud

import (
	"encoding/json"
	"net/http"
	"strconv"

	res "github.com/lautarojayat/backoffice/api/http/response"
)

func (pmux *productsMux) listProducts(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	c, err := pmux.r.List(offset, limit)

	if err != nil {
		pmux.l.Printf("couldn't retreive products list. error=%q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p := res.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	response := res.PaginatedProducts{
		Pagination: p,
		Products:   res.ToProductsResponse(c),
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		pmux.l.Printf("couldn't marshal response to json. error=%q", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
