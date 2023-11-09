package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestDecoder(t *testing.T) {
	okToken := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	expiredToken := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	privKey, pubKey := ReadKeys(t)

	tests := []struct {
		hasErr     bool
		hasToken   bool
		inputToken JWTStructure
	}{
		{false, true, JWTStructure{1, "John", "john@email.com", okToken}},
		{true, false, JWTStructure{0, "", "john@email.com", okToken}},
		{true, false, JWTStructure{1, "John", "", okToken}},
		{true, false, JWTStructure{1, "", "", okToken}},
		{true, false, JWTStructure{1, "John", "john@email.com", expiredToken}},
	}
	for i, test := range tests {
		tokenString := CreateRSASignedToken(test.inputToken, privKey)
		token, err := DecodeToken(tokenString, pubKey)
		if err != nil && !test.hasErr {
			t.Errorf("testcase %d. received an unexpected error. error=%q", i, err)
		}
		if err == nil && test.hasErr {
			t.Errorf("testcase %d. expected error but received nothing", i)
		}
		if token != nil && !test.hasToken {
			t.Errorf("testcase %d. error didn't expect a token but received one", i)
		}
		if token == nil && test.hasToken {
			t.Errorf("testcase %d. expected a token but didn't received one. error=%q", i, err)
		}

	}
}

func TestAuthChecker(t *testing.T) {
	okToken := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	expiredToken := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	privKey, pubKey := ReadKeys(t)

	tests := []struct {
		expectedCode int
		inputToken   *JWTStructure
	}{
		{http.StatusForbidden, &JWTStructure{0, "", "john@email.com", okToken}},
		{http.StatusForbidden, &JWTStructure{1, "John", "", okToken}},
		{http.StatusForbidden, &JWTStructure{1, "", "", okToken}},
		{http.StatusForbidden, &JWTStructure{1, "John", "john@email.com", expiredToken}},
		{http.StatusForbidden, nil},
		{http.StatusBadGateway, &JWTStructure{1, "John", "john@email.com", okToken}},
	}
	proxy, err := NewProxy("http://localhost:8080")

	if err != nil {
		t.Fatalf("could not start proxy for testing. error=%q", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", AuthChecker(proxy, pubKey))

	for i, test := range tests {

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)

		if err != nil {
			t.Fatalf("could not create request for testing. error=%q", err)
		}

		var tokenString string

		if test.inputToken != nil {
			tokenString = CreateRSASignedToken(*test.inputToken, privKey)
		} else {
			tokenString = ""
		}

		req.Header.Add("Auth", tokenString)

		mux.ServeHTTP(rr, req)

		if rr.Result().StatusCode != test.expectedCode {
			t.Errorf("testcase %d. expected %d status code but received %d", i, test.expectedCode, rr.Result().StatusCode)
		}

	}
}
