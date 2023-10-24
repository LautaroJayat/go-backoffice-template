package product_crud

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	products "github.com/lautarojayat/e_shop/products"
)

type productsMux struct {
	l *log.Logger
	r *products.Repo
}

func NewMux(l *log.Logger, r *products.Repo) *chi.Mux {
	pMux := &productsMux{l, r}
	m := chi.NewMux()

	m.Get("/", pMux.listProducts)
	m.Post("/", pMux.createProduct)
	m.Get("/{id}", pMux.findProduct)
	m.Put("/{id}", pMux.updateProduct)
	return m
}

func RegisterMux(topLevelMux *http.ServeMux, usersMux *chi.Mux) {
	topLevelMux.Handle("/products/", http.StripPrefix("/products", usersMux))
}
