package users_crud

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	users "github.com/lautarojayat/e_shop/users"
)

type usersMux struct {
	l *log.Logger
	r *users.Repo
}

func NewMux(l *log.Logger, r *users.Repo) *chi.Mux {
	cMux := &usersMux{l, r}
	m := chi.NewMux()

	m.Get("/", cMux.listCustomers)
	m.Post("/", cMux.createCustomer)
	m.Get("/{id}", cMux.findCustomer)
	m.Put("/{id}", cMux.updateCustomer)
	return m
}

func RegisterMux(topLevelMux *http.ServeMux, usersMux *chi.Mux) {
	topLevelMux.Handle("/users/", http.StripPrefix("/users", usersMux))
}
