package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lautarojayat/backoffice/api/http/response"
	users "github.com/lautarojayat/backoffice/users"
)

func TestListCustomerHandler(t *testing.T) {
	mux, db := setupDbAndMux(t)

	expectedUsers := []users.User{
		{Name: "user1"},
		{Name: "user2"},
		{Name: "user3"},
		{Name: "user4"},
		{Name: "user5"},
	}

	for i := 0; i < len(expectedUsers); i += 1 {
		result := db.Create(&expectedUsers[i])
		if result.Error != nil {
			t.Fatalf("could not create user to test endpoint. error=%q", result.Error)
		}
	}

	req, err := http.NewRequest("GET", "/users/", nil)

	if err != nil {
		t.Fatalf("could not create correct request to test endpoint. error=%q", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	fmt.Println(rr.Result().Body)

	if rr.Code != http.StatusOK {
		t.Errorf("status must be 200, instead got %d", rr.Code)
	}

	payload := response.PaginatedCustomers{}

	json.NewDecoder(rr.Result().Body).Decode(&payload)

	for i := 0; i < len(expectedUsers); i += 1 {
		found := false
		for _, c := range payload.Customers {
			if c.Name == expectedUsers[i].Name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("could not find user %q in payload", expectedUsers[i].Name)

		}
	}

}
