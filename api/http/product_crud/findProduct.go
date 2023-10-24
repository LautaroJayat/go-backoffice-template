package product_crud

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lautarojayat/backoffice/api/http/response"
)

func (pmux *productsMux) findProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		pmux.l.Printf("couldn't extract id from url. error=%q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := pmux.r.FindById(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(response.ToProductResponse(*p))

	if err != nil {
		pmux.l.Printf("couldn't marshal product to json response. error=%q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
