package http

import (
	"log"
	"net/http"

	"github.com/lautarojayat/e_shop/api/http/meta"
	"github.com/lautarojayat/e_shop/api/http/product_crud"
	"github.com/lautarojayat/e_shop/api/http/users_crud"
	"github.com/lautarojayat/e_shop/products"
	users "github.com/lautarojayat/e_shop/users"
)

func MakeHTTPEndpoints(l *log.Logger, usersRepo *users.Repo, productsRepo *products.Repo) *http.ServeMux {
	topLevelMux := http.NewServeMux()

	metaMux := meta.NewMux(l)
	meta.RegisterMux(topLevelMux, metaMux)

	usersMux := users_crud.NewMux(l, usersRepo)
	users_crud.RegisterMux(topLevelMux, usersMux)

	productsMux := product_crud.NewMux(l, productsRepo)
	product_crud.RegisterMux(topLevelMux, productsMux)

	return topLevelMux
}
