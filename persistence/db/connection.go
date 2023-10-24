package db

import (
	"strings"

	"github.com/lautarojayat/e_shop/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func concatOptions(cfg config.DBConfig) string {
	var builder strings.Builder
	builder.WriteString("host=")
	builder.WriteString(cfg.Host)
	builder.WriteString(" user=")
	builder.WriteString(cfg.User)
	builder.WriteString(" password=")
	builder.WriteString(cfg.Password)
	builder.WriteString(" dbname=")
	builder.WriteString(cfg.DBName)
	builder.WriteString(" port=")
	builder.WriteString(cfg.Port)
	builder.WriteString(" sslmode=")
	builder.WriteString(cfg.SSLMode)
	return builder.String()
}

func NewConnection(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := concatOptions(cfg)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil

}
