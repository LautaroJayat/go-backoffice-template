package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/lautarojayat/backoffice/api/http/response"
	"github.com/lautarojayat/backoffice/roles"
	users "github.com/lautarojayat/backoffice/users"
)

func TestFindUserHandler(t *testing.T) {
	mux, db := setupDbAndMux(t)

	expectedUser := &users.User{Name: "user"}
	result := db.Create(expectedUser)

	if result.Error != nil {
		t.Fatalf("could not create user to test endpoint. error=%q", result.Error)
	}

	tests := []struct {
		expectedStatus int
		id             uint64
		auth           roles.Role
		checkBody      bool
	}{
		{http.StatusOK, uint64(expectedUser.ID), roles.ReadUser, true},
		{http.StatusInternalServerError, uint64(expectedUser.ID + 100), roles.ReadUser, false},
		{http.StatusForbidden, uint64(expectedUser.ID), roles.CreateUser, false},
		{http.StatusForbidden, uint64(expectedUser.ID), 0, false},
	}

	for _, test := range tests {

		id := strconv.FormatUint(test.id, 10)
		var builder strings.Builder
		builder.WriteString("/users/")
		builder.WriteString(id)
		path := builder.String()

		req, err := http.NewRequest("GET", path, nil)

		if err != nil {
			t.Fatalf("could not create correct request to test endpoint. error=%q", err)
			continue
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

		var payload response.UserResponse
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
}
