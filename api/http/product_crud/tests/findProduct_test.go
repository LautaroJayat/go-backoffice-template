package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/lautarojayat/backoffice/api/http/response"
	"github.com/lautarojayat/backoffice/products"
	"github.com/lautarojayat/backoffice/roles"
)

func TestFindProductHandler(t *testing.T) {
	mux, db := setupDbAndMux(t)

	expectedProduct := &products.Product{Name: "p1", Price: 10}
	result := db.Create(expectedProduct)

	if result.Error != nil {
		t.Fatalf("could not create product to test endpoint. error=%q", result.Error)
	}

	id := uint64(expectedProduct.ID)

	tests := []struct {
		id             string
		checkBody      bool
		auth           roles.Role
		expectedStatus int
	}{
		{strconv.FormatUint(id, 10), true, roles.ReadProduct, http.StatusOK},
		{strconv.FormatUint(id+1, 10), false, roles.ReadProduct, http.StatusNotFound},
		{strconv.FormatUint(id, 10), false, roles.CreateProduct, http.StatusForbidden},
	}

	for _, test := range tests {

		var builder strings.Builder
		builder.WriteString("/products/")
		builder.WriteString(test.id)
		path := builder.String()

		req, err := http.NewRequest("GET", path, nil)

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
}
