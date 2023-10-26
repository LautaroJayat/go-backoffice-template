package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/lautarojayat/backoffice/api/http/response"
	"github.com/lautarojayat/backoffice/roles"
	users "github.com/lautarojayat/backoffice/users"
)

func TestListUserHandler(t *testing.T) {
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

	tests := []struct {
		expectedStatus int
		auth           roles.Role
		checkBody      bool
	}{
		{http.StatusOK, roles.ReadUser, true},
		{http.StatusForbidden, roles.CreateProduct, false},
		{http.StatusForbidden, 0, false},
	}
	for _, test := range tests {

		req, err := http.NewRequest("GET", "/users/", nil)

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

		payload := response.PaginatedUsers{}

		json.NewDecoder(rr.Result().Body).Decode(&payload)

		for i := 0; i < len(expectedUsers); i += 1 {
			found := false
			for _, c := range payload.Users {
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

}
