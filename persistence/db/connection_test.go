package db

import (
	"testing"

	"github.com/lautarojayat/backoffice/config"
)

func TestConnectionOptions(t *testing.T) {
	configs := config.DBConfig{
		Host:     "host",
		User:     "user",
		Password: "password",
		DBName:   "dbname",
		Port:     "port",
		SSLMode:  "sslmode",
	}
	expectedString := "host=host user=user password=password dbname=dbname port=port sslmode=sslmode"
	o := concatOptions(configs)
	if o != expectedString {
		t.Errorf("bad connection options. expected=%q got=%q", expectedString, o)
	}
}
