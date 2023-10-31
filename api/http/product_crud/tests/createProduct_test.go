package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/lautarojayat/backoffice/api/http/response"
	"github.com/lautarojayat/backoffice/roles"
)

func TestCreateCustomerHandler(t *testing.T) {

	tests := []struct {
		body           string
		checkBody      bool
		auth           roles.Role
		expectedStatus int
	}{
		{"{\"name\":\"product1\",\"price\":10}", true, roles.CreateProduct, http.StatusCreated},
		{"{\"name\":\"product1\",\"price\":10}", false, roles.CreateUser, http.StatusForbidden},
		{"{\"nam\":\"product1\",\"price\":10}", false, roles.CreateProduct, http.StatusBadRequest},
		{"{\"name\":\"product1\",\"rice\":10}", false, roles.CreateProduct, http.StatusBadRequest},
		{"{\"name\":\"\",\"price\":\"\"}", false, roles.CreateProduct, http.StatusBadRequest},
	}

	mux, _ := setupDbAndMux(t)

	for _, test := range tests {

		var buf bytes.Buffer
		buf.WriteString(test.body)

		req, err := http.NewRequest("POST", "/products/", bytes.NewReader(buf.Bytes()))

		if err != nil {
			t.Fatalf("could not create correct request to test endpoint. error=%q", err)
		}

		req.Header.Add(roles.DecodedPermsHeader, strconv.FormatUint(uint64(test.auth), 2))

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		status := rr.Code
		if status != test.expectedStatus {
			t.Errorf("status must be %d, instead got %d", test.expectedStatus, status)
		}

		if !test.checkBody {
			continue
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
}
