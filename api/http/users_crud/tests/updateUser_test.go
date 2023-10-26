package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/lautarojayat/backoffice/roles"
	users "github.com/lautarojayat/backoffice/users"
)

func TestUpdateCustomerHandler(t *testing.T) {
	mux, db := setupDbAndMux(t)

	initialUser := &users.User{Name: "user"}
	result := db.Create(initialUser)

	if result.Error != nil {
		t.Fatalf("could not create user to test endpoint. error=%q", result.Error)
	}

	tests := []struct {
		expectedStatus int
		id             uint64
		jsonStructure  string
		newName        string
		auth           roles.Role
		checkDB        bool
	}{
		{http.StatusAccepted, uint64(initialUser.ID), `{"name":"%s"}`, "user2", roles.ModifyUser, true},
		{http.StatusBadRequest, uint64(initialUser.ID), `{"nam":"%s"}`, "user1", roles.ModifyUser, false},
		{http.StatusBadRequest, uint64(initialUser.ID), "%s", "", roles.ModifyUser, false},
		{http.StatusForbidden, uint64(initialUser.ID), `{"name":"%s"}`, "", roles.ReadUser, false},
		{http.StatusForbidden, uint64(initialUser.ID), `{"name":"%s"}`, "user5", 0, false},
		{http.StatusNotFound, uint64(initialUser.ID + 600), `{"name":"%s"}`, "user6", roles.ModifyUser, false},
	}

	for _, test := range tests {

		id := strconv.FormatUint(test.id, 10)

		var builder strings.Builder
		builder.WriteString("/users/")
		builder.WriteString(id)
		path := builder.String()

		reqBody := fmt.Sprintf(test.jsonStructure, test.newName)
		req, err := http.NewRequest("PUT", path, bytes.NewReader([]byte(reqBody)))

		if err != nil {
			t.Fatalf("could not create correct request to test endpoint. error=%q", err)
		}

		req.Header.Add(roles.DecodedPermsHeader, strconv.FormatUint(uint64(test.auth), 2))

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != test.expectedStatus {
			t.Errorf("status must be %d, instead got %d", test.expectedStatus, rr.Code)
		}
		if !test.checkDB {
			continue
		}

		updatedUser := users.User{}

		result = db.Where("id = ?", test.id).Find(&updatedUser)

		if result.Error != nil {
			t.Fatalf("couldn't read updated user from db. error=%q", result.Error)
		}

		if updatedUser.ID != uint(test.id) {
			t.Errorf("Id must be %d, instead got %d", test.id, updatedUser.ID)
		}

		if updatedUser.Name != test.newName {
			t.Errorf("expected name %q, but got %q", test.newName, updatedUser.Name)
		}

		if updatedUser.CreatedAt.Local().UnixMilli() != initialUser.CreatedAt.UnixMilli() {
			t.Errorf("CreatedAt must be %d, instead got %d", initialUser.CreatedAt.UnixMilli(), updatedUser.CreatedAt.UnixMilli())
		}
		if updatedUser.UpdatedAt.UnixMilli() <= initialUser.UpdatedAt.UnixMilli() {
			t.Errorf("UpdatedAt must be greater after an update, updated=%d before=%d", updatedUser.UpdatedAt.UnixMilli(), initialUser.UpdatedAt.UnixMilli())
		}
	}
}
