package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/lautarojayat/backoffice/products"
)

func TestUpdateCustomerHandler(t *testing.T) {
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

	newName := "p2"
	newPrice := 15
	reqBody := fmt.Sprintf(`{"name":"%s", "price":%d}`, newName, newPrice)
	req, err := http.NewRequest("PUT", path, bytes.NewReader([]byte(reqBody)))

	if err != nil {
		t.Fatalf("could not create correct request to test endpoint. error=%q", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("status must be 202, instead got %d", rr.Code)
	}

	updatedProduct := products.Product{}

	result = db.Where("id = ?", expectedProduct.ID).Find(&updatedProduct)

	if result.Error != nil {
		t.Fatalf("couldn't read updated user from db. error=%q", result.Error)
	}

	if updatedProduct.ID != expectedProduct.ID {
		t.Errorf("Id must be %d, instead got %d", expectedProduct.ID, updatedProduct.ID)
	}

	if updatedProduct.Name != newName {
		t.Errorf("expected name must be %q, but got %q", newName, updatedProduct.Name)
	}
	if int(updatedProduct.Price) != newPrice {
		t.Errorf("expected price must be %d, but got %d", newPrice, updatedProduct.Price)
	}

	if updatedProduct.CreatedAt.Local().UnixMilli() != expectedProduct.CreatedAt.UnixMilli() {
		t.Errorf("CreatedAt must be %d, instead got %d", expectedProduct.CreatedAt.UnixMilli(), updatedProduct.CreatedAt.UnixMilli())
	}
	if updatedProduct.UpdatedAt.UnixMilli() <= expectedProduct.UpdatedAt.UnixMilli() {
		t.Errorf("UpdatedAt must be greater after an update")
	}
}
