package product_crud

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/lautarojayat/backoffice/api/http/middleware"
	products "github.com/lautarojayat/backoffice/products"
	"github.com/lautarojayat/backoffice/roles"
)

type productsMux struct {
	l *log.Logger
	r *products.Repo
}

func NewMux(l *log.Logger, r *products.Repo) *chi.Mux {
	pMux := &productsMux{l, r}
	m := chi.NewMux()

	m.Get("/", middleware.CheckRole(l, []roles.Role{roles.ReadProduct}, pMux.listProducts))
	m.Post("/", middleware.CheckRole(l, []roles.Role{roles.CreateProduct}, pMux.createProduct))
	m.Get("/{id}", middleware.CheckRole(l, []roles.Role{roles.ReadProduct}, pMux.findProduct))
	m.Put("/{id}", middleware.CheckRole(l, []roles.Role{roles.ModifyProduct}, pMux.updateProduct))
	return m
}

func RegisterMux(topLevelMux *http.ServeMux, usersMux *chi.Mux) {
	topLevelMux.Handle("/products/", http.StripPrefix("/products", usersMux))
}
