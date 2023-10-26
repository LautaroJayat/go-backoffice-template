package users_crud

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	users "github.com/lautarojayat/backoffice/users"
)

func (cmux *usersMux) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		cmux.l.Printf("couldn't extract id from url. error=%q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	users := &users.User{}

	if err := json.NewDecoder(r.Body).Decode(users); err != nil {
		cmux.l.Printf("could not extract user from request. error=%q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if users.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	done, err := cmux.r.UpdateOne(id, users.Name)

	if err != nil {
		cmux.l.Printf("could not update users. error=%q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !done {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
