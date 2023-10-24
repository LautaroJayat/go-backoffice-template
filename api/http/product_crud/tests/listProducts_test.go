package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lautarojayat/e_shop/api/http/response"
	"github.com/lautarojayat/e_shop/products"
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

	for i := 0; i < len(expectedProducts); i += 1 {
		result := db.Create(&expectedProducts[i])
		if result.Error != nil {
			t.Fatalf("could not create product to test endpoint. error=%q", result.Error)
		}
	}

	req, err := http.NewRequest("GET", "/products/", nil)

	if err != nil {
		t.Fatalf("could not create correct request to test endpoint. error=%q", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	fmt.Println(rr.Result().Body)

	if rr.Code != http.StatusOK {
		t.Errorf("status must be 200, instead got %d", rr.Code)
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
