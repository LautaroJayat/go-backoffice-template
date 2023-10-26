package users_crud

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/lautarojayat/backoffice/api/http/middleware"
	"github.com/lautarojayat/backoffice/roles"
	users "github.com/lautarojayat/backoffice/users"
)

type usersMux struct {
	l *log.Logger
	r *users.Repo
}

func NewMux(l *log.Logger, r *users.Repo) *chi.Mux {
	uMux := &usersMux{l, r}
	m := chi.NewMux()

	m.Get("/", middleware.CheckRole(l, []roles.Role{roles.ReadUser}, uMux.listUsers))
	m.Post("/", middleware.CheckRole(l, []roles.Role{roles.CreateUser}, uMux.createUser))
	m.Get("/{id}", middleware.CheckRole(l, []roles.Role{roles.ReadUser}, uMux.findUser))
	m.Put("/{id}", middleware.CheckRole(l, []roles.Role{roles.ModifyUser}, uMux.updateUser))
	return m
}

func RegisterMux(topLevelMux *http.ServeMux, usersMux *chi.Mux) {
	topLevelMux.Handle("/users/", http.StripPrefix("/users", usersMux))
}
