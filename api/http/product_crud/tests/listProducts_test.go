package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/lautarojayat/backoffice/api/http/response"
	"github.com/lautarojayat/backoffice/products"
	"github.com/lautarojayat/backoffice/roles"
)

func TestListCustomerHandler(t *testing.T) {
	mux, db := setupDbAndMux(t)

	expectedProducts := []products.Product{
		{Name: "p1", Price: 10},
		{Name: "p2", Price: 10},
		{Name: "p3", Price: 10},
		{Name: "p4", Price: 10},
		{Name: "p5", Price: 10},
	}

	tests := []struct {
		auth           roles.Role
		expectedStatus int
		checkBody      bool
	}{
		{roles.ReadProduct, http.StatusOK, true},
		{roles.CreateProduct, http.StatusForbidden, false},
	}

	for i := 0; i < len(expectedProducts); i += 1 {
		result := db.Create(&expectedProducts[i])
		if result.Error != nil {
			t.Fatalf("could not create product to test endpoint. error=%q", result.Error)
		}
	}

	for _, test := range tests {

		req, err := http.NewRequest("GET", "/products/", nil)

		if err != nil {
			t.Fatalf("could not create correct request to test endpoint. error=%q", err)
		}

		req.Header.Add(roles.DecodedPermsHeader, strconv.FormatUint(uint64(test.auth), 2))

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		fmt.Println(rr.Result().Body)

		if rr.Code != test.expectedStatus {
			t.Errorf("status must be %d, instead got %d", test.expectedStatus, rr.Code)
		}

		if !test.checkBody {
			continue
		}

		payload := response.PaginatedProducts{}

		json.NewDecoder(rr.Result().Body).Decode(&payload)

		for i := 0; i < len(expectedProducts); i += 1 {
			found := false
			for _, c := range payload.Products {
				if c.Name == expectedProducts[i].Name {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("could not find user %q in payload", expectedProducts[i].Name)

			}
		}
	}

}
