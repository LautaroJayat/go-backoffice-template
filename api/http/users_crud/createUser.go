package users_crud

import (
	"encoding/json"
	"net/http"

	"github.com/lautarojayat/backoffice/api/http/response"
	users "github.com/lautarojayat/backoffice/users"
)

func (umux *usersMux) createUser(w http.ResponseWriter, r *http.Request) {
	users := &users.User{}

	if err := json.NewDecoder(r.Body).Decode(users); err != nil {
		umux.l.Printf("could not extract users from request. error=%q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(users.Name) < 4 {
		umux.l.Println("user name should contian at least 4 chars")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := umux.r.CreateOne(users.Name)

	if err != nil {
		umux.l.Printf("could not create users. error=%q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(response.ToUserResponse(*users))
	if err != nil {
		umux.l.Printf("could not format json for response. error=%q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(responseBody)
}
