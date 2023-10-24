package meta

import (
	"fmt"
	"net/http"
)

func (m *metaMux) ready(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("ok"))

	if err != nil {
		m.logError(fmt.Sprintf("something happened while reporting readyness: %q", err))
	}

}

func (m *metaMux) health(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("ok"))

	if err != nil {
		m.logError(fmt.Sprintf("something happened while reporting status: %q", err))
	}

}
