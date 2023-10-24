package product_crud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lautarojayat/e_shop/products"
)

func (pmux *productsMux) updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		pmux.l.Printf("couldn't extract id from url. error=%q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p := products.Product{}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		pmux.l.Printf("could not extract product from request. error=%q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if p.Name == "" && p.Price == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"must provide name or price"}`)
		return
	}

	err = pmux.r.UpdateOne(id, p)

	if err != nil {
		pmux.l.Printf("could not update users. error=%q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
