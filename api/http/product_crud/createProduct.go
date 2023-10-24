package product_crud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lautarojayat/e_shop/api/http/response"
	"github.com/lautarojayat/e_shop/products"
)

func (pmux *productsMux) createProduct(w http.ResponseWriter, r *http.Request) {
	p := &products.Product{}

	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		pmux.l.Printf("could not extract product from request. error=%q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	errors := make([]string, 0, 2)
	if p.Name == "" {
		errors = append(errors, `"missing name in request"`)
	}

	if p.Price == 0 {
		errors = append(errors, `"missing price in request"`)
	}

	if len(errors) > 0 {
		e := strings.Join(errors, ",")
		pmux.l.Printf("could not create product due to bad request. error=%q", e)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"errors":[%s]"}`, e)
		return

	}

	p, err := pmux.r.CreateOne(p.Name, p.Price)

	if err != nil {
		pmux.l.Printf("could not create product. error=%q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(response.ToProductResponse(*p))

	if err != nil {
		pmux.l.Printf("could not format json for response. error=%q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(responseBody)
}
