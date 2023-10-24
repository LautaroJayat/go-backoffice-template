package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lautarojayat/e_shop/api/http/response"
)

func TestCreateCustomerHandler(t *testing.T) {

	mux, _ := setupDbAndMux(t)

	var buf bytes.Buffer
	buf.WriteString("{\"name\":\"product1\",\"price\":10}")

	req, err := http.NewRequest("POST", "/products/", bytes.NewReader(buf.Bytes()))

	if err != nil {
		t.Fatalf("could not create correct request to test endpoint. error=%q", err)
	}

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusCreated {
		t.Errorf("status must be 201, instead got %d", status)
	}

	payload := response.ProductResponse{}

	err = json.NewDecoder(rr.Result().Body).Decode(&payload)

	if err != nil {
		t.Fatalf("couldn't read response body. error=%q", err)
	}

	if payload.Id <= 0 {
		t.Errorf("Id must have been greater than 0, instead got %d", payload.Id)
	}
	if payload.CreatedAt <= 0 {
		t.Errorf("CreatedAt must have been greater than 0, instead got %d", payload.CreatedAt)
	}
	if payload.UpdatedAt <= 0 {
		t.Errorf("UpdatedAt must have been greater than 0, instead got %d", payload.UpdatedAt)
	}
	if payload.Name != "product1" {
		t.Errorf("product name must have been product1, instead got %q", payload.Name)
	}
	if payload.Price != 10 {
		t.Errorf("price must have been 10, instead got %d", payload.Price)
	}
}
