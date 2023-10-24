package meta

import (
	"log"
	"net/http"
)

type metaMux struct {
	l *log.Logger
}

func (m *metaMux) logError(e string) {
	m.l.Printf("error=%q", e)
}

func NewMux(logger *log.Logger) *http.ServeMux {
	metaMux := &metaMux{logger}
	m := http.NewServeMux()
	m.HandleFunc("/ready", metaMux.ready)
	m.HandleFunc("/status", metaMux.health)
	return m
}

func RegisterMux(topLevelMux *http.ServeMux, metaMux *http.ServeMux) {
	topLevelMux.Handle("/meta/", http.StripPrefix("/meta", metaMux))
}
