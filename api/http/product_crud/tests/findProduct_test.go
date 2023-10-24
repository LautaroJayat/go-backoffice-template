package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/lautarojayat/e_shop/api/http/response"
	"github.com/lautarojayat/e_shop/products"
)

func TestFindProductHandler(t *testing.T) {
	mux, db := setupDbAndMux(t)

	expectedProduct := &products.Product{Name: "p1", Price: 10}
	result := db.Create(expectedProduct)

	if result.Error != nil {
		t.Fatalf("could not create product to test endpoint. error=%q", result.Error)
	}

	id := strconv.FormatUint(uint64(expectedProduct.ID), 10)

	var builder strings.Builder
	builder.WriteString("/products/")
	builder.WriteString(id)
	path := builder.String()

	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		t.Fatalf("could not create correct request to test endpoint. error=%q", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status must be 200, instead got %d", rr.Code)
	}

	var payload response.ProductResponse
	err = json.NewDecoder(rr.Result().Body).Decode(&payload)

	if err != nil {
		t.Fatalf("couldn't read response body. error=%q", err)
	}

	if payload.Id != expectedProduct.ID {
		t.Errorf("Id must be %d, instead got %d", expectedProduct.ID, payload.Id)
	}
	if payload.Name != expectedProduct.Name {
		t.Errorf("Name must be %q, instead got %q", expectedProduct.Name, payload.Name)
	}
	if payload.Price != expectedProduct.Price {
		t.Errorf("Price must be %d, instead got %d", expectedProduct.Price, payload.Price)
	}
	if payload.CreatedAt != expectedProduct.CreatedAt.UnixMilli() {
		t.Errorf("CreatedAt must be %d, instead got %d", expectedProduct.CreatedAt.UnixMilli(), payload.CreatedAt)
	}
	if payload.UpdatedAt != expectedProduct.CreatedAt.UnixMilli() {
		t.Errorf("UpdatedAt must be %d, instead got %d", expectedProduct.CreatedAt.UnixMilli(), payload.UpdatedAt)
	}
}
