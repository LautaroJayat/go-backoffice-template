package proxy

import (
	"crypto/rsa"
	"io"
	"log"
	"os"
	"path"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

// This method creates a signed token.
// As this proxy is not an auth provider. This method is just to create inputs for testing.
func CreateRSASignedToken(input JWTStructure, privKey *rsa.PrivateKey) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, input)
	tokenString, err := token.SignedString(privKey)
	if err != nil {
		log.Fatalf("could not create jwt for testing. error=%q", err)
	}
	return tokenString
}

// This method reads the keys generated by make gen-test-keys.
// Is used only for testing
func ReadKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
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
