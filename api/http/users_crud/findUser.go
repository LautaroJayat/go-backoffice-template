package users_crud

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lautarojayat/backoffice/api/http/response"
)

func (cmux *usersMux) findUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		cmux.l.Printf("couldn't extract id from url. error=%q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c, err := cmux.r.FindById(id)
	if err != nil {
		cmux.l.Printf("internal error. error=%q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if c == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(response.ToUserResponse(*c))

	if err != nil {
		cmux.l.Printf("couldn't marshal users to json response. error=%q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
