package meta

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadyHandler(t *testing.T) {

	metaMux := NewMux(log.Default())

	gralMux := http.NewServeMux()
	gralMux.Handle("/meta/", http.StripPrefix("/meta", metaMux))

	req, err := http.NewRequest("GET", "/meta/ready", nil)
	if err != nil {
		t.Fatalf("could not create correct request to test readyness endpoint. error=%q", err)
	}

	rr := httptest.NewRecorder()

	gralMux.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("status must be 200, instead got %d", status)
	}

	body := make([]byte, 2)
	rr.Body.Read(body)
	bodyContent := string(body)
	if bodyContent != "ok" {
		t.Errorf("body must be 'ok', instead got %q", bodyContent)
	}

}

func TestStatusHandler(t *testing.T) {

	metaMux := NewMux(log.Default())

	gralMux := http.NewServeMux()
	gralMux.Handle("/meta/", http.StripPrefix("/meta", metaMux))

	req, err := http.NewRequest("GET", "/meta/status", nil)
	if err != nil {
		t.Fatalf("could not create correct request to test status endpoint. error=%q", err)
	}

	rr := httptest.NewRecorder()

	gralMux.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("status must be 200, instead got %d", status)
	}

	body := make([]byte, 2)
	rr.Body.Read(body)
	bodyContent := string(body)
	if bodyContent != "ok" {
		t.Errorf("body must be 'ok', instead got %q", bodyContent)
	}

}
