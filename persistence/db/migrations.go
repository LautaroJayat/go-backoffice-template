package db

import (
	"log"

	"github.com/lautarojayat/e_shop/products"
	users "github.com/lautarojayat/e_shop/users"
	"gorm.io/gorm"
)

func RunMigrations(l *log.Logger, db *gorm.DB) error {
	err := db.AutoMigrate(&users.User{})
	if err != nil {
		return err
	}
	l.Println("User model migrated OK")

	err = db.AutoMigrate(&products.Product{})
	if err != nil {
		return err
	}
	l.Println("Product model migrated OK")
	return nil
}
