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

func TestCreateUserHandler(t *testing.T) {

	tests := []struct {
		expectedStatus int
		body           string
		auth           roles.Role
		checkBody      bool
	}{
		{http.StatusCreated, `{"name": "user"}`, roles.CreateUser, true},
		{http.StatusBadRequest, "", roles.CreateUser, false},
		{http.StatusBadRequest, `{"nam": "user"}`, roles.CreateUser, false},
		{http.StatusForbidden, `{"name": "user"}`, roles.ReadUser, false},
		{http.StatusForbidden, `{"name": "user"}`, 0, false},
	}

	mux, _ := setupDbAndMux(t)

	for index, test := range tests {

		var buf bytes.Buffer
		buf.WriteString(test.body)

		req, err := http.NewRequest("POST", "/users/", bytes.NewReader(buf.Bytes()))

		if err != nil {
			t.Fatalf("could not create correct request to test endpoint. error=%q", err)
		}
		req.Header.Add(roles.DecodedPermsHeader, strconv.FormatUint(uint64(test.auth), 2))

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != test.expectedStatus {
			t.Errorf("test number %d: status must be %d, instead got %d", index, test.expectedStatus, rr.Code)
		}

		if !test.checkBody {
			continue
		}

		var payload response.UserResponse

		err = json.NewDecoder(rr.Result().Body).Decode(&payload)

		if err != nil {
			t.Errorf("couldn't read response body. error=%q", err)
			continue
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
	}
}
