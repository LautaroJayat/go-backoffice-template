package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	ro "github.com/lautarojayat/backoffice/roles"
)

func dummyEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestCheckRole(t *testing.T) {
	tests := []struct {
		guardedWith    []ro.Role
		inputRole      ro.Role
		expectedHeader int
	}{
		{
			[]ro.Role{ro.ReadUser},
			ro.ReadUser,
			http.StatusOK,
		},
		{
			[]ro.Role{ro.ReadUser},
			ro.CreateUser,
			http.StatusForbidden,
		},
		{
			[]ro.Role{ro.ReadUser},
			ro.SuperAdmin,
			http.StatusOK,
		},
		{
			[]ro.Role{ro.ReadUser},
			ro.UserAdmin,
			http.StatusOK,
		},
		{
			[]ro.Role{ro.ReadUser, ro.CreateUser},
			ro.ReadUser,
			http.StatusForbidden,
		},
		{
			[]ro.Role{ro.ReadUser, ro.CreateUser, ro.ModifyUser},
			ro.UserAdmin,
			http.StatusOK,
		},
		{
			[]ro.Role{},
			ro.ReadUser,
			http.StatusOK,
		},
	}
	l := log.Default()

	for _, test := range tests {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Error("couldn't prepare request for testing middleware")
			t.FailNow()
		}
		req.Header.Add(ro.DecodedPermsHeader, strconv.FormatUint(uint64(test.inputRole), 2))
		CheckRole(l, test.guardedWith, dummyEndpoint).ServeHTTP(rr, req)
		status := rr.Result().StatusCode
		if status != test.expectedHeader {
			t.Errorf("expected %d status when guard was %#v, and provided role was %d. got=%d",
				test.expectedHeader,
				test.guardedWith,
				test.inputRole,
				status)
		}

	}
}
