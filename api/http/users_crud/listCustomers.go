package users_crud

import (
	"encoding/json"
	"net/http"
	"strconv"

	res "github.com/lautarojayat/backoffice/api/http/response"
)

func (cmux *usersMux) listCustomers(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	c, err := cmux.r.List(offset, limit)

	if err != nil {
		cmux.l.Printf("couldn't retreive users list. error=%q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p := res.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	response := &res.PaginatedCustomers{
		Pagination: p,
		Customers:  res.ToCustomersResponse(c),
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		cmux.l.Printf("couldn't marshal response to json. error=%q", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
