package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	users "github.com/lautarojayat/e_shop/users"
)

func TestUpdateCustomerHandler(t *testing.T) {
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

	newName := "user2"
	reqBody := fmt.Sprintf(`{"name":"%s"}`, newName)
	req, err := http.NewRequest("PUT", path, bytes.NewReader([]byte(reqBody)))

	if err != nil {
		t.Fatalf("could not create correct request to test endpoint. error=%q", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("status must be 202, instead got %d", rr.Code)
	}

	updatedUser := users.User{}

	result = db.Where("id = ?", expectedUser.ID).Find(&updatedUser)

	if result.Error != nil {
		t.Fatalf("couldn't read updated user from db. error=%q", result.Error)
	}

	if updatedUser.ID != expectedUser.ID {
		t.Errorf("Id must be %d, instead got %d", expectedUser.ID, updatedUser.ID)
	}

	if updatedUser.Name != newName {
		t.Errorf("expected name %q, but got %q", newName, updatedUser.Name)
	}

	if updatedUser.CreatedAt.Local().UnixMilli() != expectedUser.CreatedAt.UnixMilli() {
		t.Errorf("CreatedAt must be %d, instead got %d", expectedUser.CreatedAt.UnixMilli(), updatedUser.CreatedAt.UnixMilli())
	}
	if updatedUser.UpdatedAt.UnixMilli() <= expectedUser.UpdatedAt.UnixMilli() {
		t.Errorf("UpdatedAt must be greater after an update")
	}
}
