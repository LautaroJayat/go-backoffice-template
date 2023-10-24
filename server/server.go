package server

import (
	"net/http"
	"time"

	"github.com/lautarojayat/e_shop/config"
)

func NewServer(cfg config.HTTPConfig, mux *http.ServeMux) *http.Server {
	s := &http.Server{
		Addr:           cfg.Port,
		Handler:        mux,
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}
	return s
}
