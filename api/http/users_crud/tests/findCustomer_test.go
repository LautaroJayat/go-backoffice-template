package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/lautarojayat/backoffice/api/http/response"
	users "github.com/lautarojayat/backoffice/users"
)

func TestFindCustomerHandler(t *testing.T) {
	mux, db := setupDbAndMux(t)

	expectedUser := &users.User{Name: "user"}
	result := db.Create(expectedUser)

	if result.Error != nil {
		t.Fatalf("could not create user to test endpoint. error=%q", result.Error)
	}

	id := strconv.FormatUint(uint64(expectedUser.ID), 10)

	var builder strings.Builder
	builder.WriteString("/users/")
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

	var payload response.CustomerResponse
	err = json.NewDecoder(rr.Result().Body).Decode(&payload)

	if err != nil {
		t.Fatalf("couldn't read response body. error=%q", err)
	}

	if payload.Id != expectedUser.ID {
		t.Errorf("Id must be %d, instead got %d", expectedUser.ID, payload.Id)
	}
	if payload.CreatedAt != expectedUser.CreatedAt.UnixMilli() {
		t.Errorf("CreatedAt must be %d, instead got %d", expectedUser.CreatedAt.UnixMilli(), payload.CreatedAt)
	}
	if payload.UpdatedAt != expectedUser.CreatedAt.UnixMilli() {
		t.Errorf("UpdatedAt must be %d, instead got %d", expectedUser.CreatedAt.UnixMilli(), payload.UpdatedAt)
	}
}
