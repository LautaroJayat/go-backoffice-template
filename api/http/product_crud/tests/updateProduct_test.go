package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/lautarojayat/backoffice/products"
	"github.com/lautarojayat/backoffice/roles"
)

func TestUpdateCustomerHandler(t *testing.T) {
	mux, db := setupDbAndMux(t)

	expectedProduct := &products.Product{Name: "p1", Price: 10}
	result := db.Create(expectedProduct)

	if result.Error != nil {
		t.Fatalf("could not create product to test endpoint. error=%q", result.Error)
	}

	id := uint64(expectedProduct.ID)

	tests := []struct {
		id             string
		body           string
		checkBody      bool
		auth           roles.Role
		expectedStatus int
		expectedName   string
		expectedPrice  int
	}{
		{strconv.FormatUint(id, 10), "{\"name\":\"product2\",\"price\":11}", true, roles.ModifyProduct, http.StatusAccepted, "product2", 11},
		{strconv.FormatUint(id+1, 10), "{\"name\":\"product1\",\"price\":10}", false, roles.ModifyProduct, http.StatusNotFound, "", 0},
		{strconv.FormatUint(id, 10), "{\"nam\":\"product3\",\"price\":12}", true, roles.ModifyProduct, http.StatusAccepted, "product2", 12},
		{strconv.FormatUint(id, 10), "{\"name\":\"product4\",\"rice\":10}", true, roles.ModifyProduct, http.StatusAccepted, "product4", 12},
		{strconv.FormatUint(id, 10), "{\"name\":\"\",\"price\":\"\"}", false, roles.ModifyProduct, http.StatusBadRequest, "", 0},
	}

	for _, test := range tests {

		var builder strings.Builder
		builder.WriteString("/products/")
		builder.WriteString(test.id)
		path := builder.String()

		req, err := http.NewRequest("PUT", path, bytes.NewReader([]byte(test.body)))

		if err != nil {
			t.Fatalf("could not create correct request to test endpoint. error=%q", err)
		}

		req.Header.Add(roles.DecodedPermsHeader, strconv.FormatUint(uint64(test.auth), 2))

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != test.expectedStatus {
			t.Errorf("status must be %d, instead got %d", test.expectedStatus, rr.Code)
		}

		if !test.checkBody {
			continue
		}

		id, err := strconv.Atoi(test.id)
		if err != nil {
			t.Error("could not convert string id to int for query")
		}

		updatedProduct := products.Product{}

		result = db.Where("id = ?", id).Find(&updatedProduct)

		if result.Error != nil {
			t.Fatalf("couldn't read updated product from db. error=%q", result.Error)
		}

		if updatedProduct.ID != uint(id) {
			t.Errorf("Id must be %d, instead got %d", id, updatedProduct.ID)
		}

		if updatedProduct.Name != test.expectedName {
			t.Errorf("expected name must be %q, but got %q", test.expectedName, updatedProduct.Name)
		}
		if int(updatedProduct.Price) != test.expectedPrice {
			t.Errorf("expected price must be %d, but got %d", test.expectedPrice, updatedProduct.Price)
		}

		if updatedProduct.CreatedAt.Local().UnixMilli() != expectedProduct.CreatedAt.UnixMilli() {
			t.Errorf("CreatedAt must be %d, instead got %d", expectedProduct.CreatedAt.UnixMilli(), updatedProduct.CreatedAt.UnixMilli())
		}
		if updatedProduct.UpdatedAt.UnixMilli() <= expectedProduct.UpdatedAt.UnixMilli() {
			t.Errorf("UpdatedAt must be greater after an update")
		}
	}
}
