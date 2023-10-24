package tests

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/lautarojayat/e_shop/api/http/users_crud"
	"github.com/lautarojayat/e_shop/config"
	database "github.com/lautarojayat/e_shop/persistence/db"
	users "github.com/lautarojayat/e_shop/users"
	"gorm.io/gorm"
)

func MuckPublisher(ctx context.Context, payload users.UsersOp) {}

func cleanup(db *gorm.DB) {
	db.Unscoped().Where("id > ?", 0).Delete(users.User{})
}

func setupMux(db *gorm.DB) *http.ServeMux {
	topLevelMux := http.NewServeMux()
	repo := users.NewRepo(db, MuckPublisher)
	usersMux := users_crud.NewMux(log.Default(), repo)
	users_crud.RegisterMux(topLevelMux, usersMux)
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
