package proxy

import (
	"crypto/rsa"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createRSASignedToken(input JWTStructure, privKey *rsa.PrivateKey) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, input)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(privKey)
	if err != nil {
		log.Fatalf("could not create jwt for testing. error=%q", err)
	}
	return tokenString
}

func readKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	cwd, err := os.Getwd()

	if err != nil {
		t.Fatalf("could not acces cwd. error=%q", err)
	}

	privateKeyPath := path.Join(cwd, ".tmp", "private.pem")
	publicKeyPath := path.Join(cwd, ".tmp", "public.pem")

	privFile, err := os.Open(privateKeyPath)

	if err != nil {
		t.Fatalf("could not open private key. error=%q", err)
	}

	pubFile, err := os.Open(publicKeyPath)

	if err != nil {
		t.Fatalf("could not open public key. error=%q", err)
	}

	privateKey, err := io.ReadAll(privFile)

	if err != nil {
		t.Fatalf("could not read private key. error=%q", err)
	}

	pubKey, err := io.ReadAll(pubFile)

	if err != nil {
		t.Fatalf("could not read public key. error=%q", err)
	}

	priv, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)

	if err != nil {
		t.Fatalf("could not parse rsa private key. error=%q", err)
	}

	pub, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)

	if err != nil {
		t.Fatalf("could not parse rsa private key. error=%q", err)
	}

	return priv, pub

}

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
	privKey, pubKey := readKeys(t)

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
		tokenString := createRSASignedToken(test.inputToken, privKey)
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

func TestAuthchecker(t *testing.T) {
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
	privKey, pubKey := readKeys(t)

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
			tokenString = createRSASignedToken(*test.inputToken, privKey)
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
