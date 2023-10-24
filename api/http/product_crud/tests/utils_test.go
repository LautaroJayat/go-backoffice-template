package tests

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/lautarojayat/e_shop/api/http/product_crud"
	"github.com/lautarojayat/e_shop/config"
	database "github.com/lautarojayat/e_shop/persistence/db"
	"github.com/lautarojayat/e_shop/products"
	"gorm.io/gorm"
)

func cleanup(db *gorm.DB) {
	db.Unscoped().Where("id > ?", 0).Delete(products.Product{})
}

func FakePublisher(context.Context, products.ProductOp) {}

func setupMux(db *gorm.DB) *http.ServeMux {
	topLevelMux := http.NewServeMux()
	repo := products.NewRepo(db, FakePublisher)
	productsMux := product_crud.NewMux(log.Default(), repo)
	product_crud.RegisterMux(topLevelMux, productsMux)
	return topLevelMux

}

func setupDbAndMux(t *testing.T) (*http.ServeMux, *gorm.DB) {
	db, err := database.NewConnection(config.DBConfig{
		Host:     "localhost",
		User:     "user",
		Password: "password",
		DBName:   "db1",
		Port:     "5432",
		SSLMode:  "disable"})

	if err != nil {
		t.Fatalf("could not connect with the db. error=%q", err)
	}

	database.RunMigrations(log.Default(), db)

	cleanup(db)

	mux := setupMux(db)
	return mux, db
}
