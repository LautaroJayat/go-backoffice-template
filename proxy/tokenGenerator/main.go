package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lautarojayat/backoffice/proxy"
)

func main() {
	keyPath := os.Args[1]
	overrideExp := os.Args[2] == "update"
	if keyPath == "" {
		panic("please provide path to private key as first command line argument")
	}
	privateKey, _ := proxy.ReadKeys(keyPath)
	path := "./jwt-sample.json"
	tokenBytes, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("could not read json sample. error=%q", err))
	}

	claims := proxy.JWTStructure{}
	err = json.Unmarshal(tokenBytes, &claims)
	if err != nil {
		panic(fmt.Sprintf("could not unmarshall token into struct. error=%q", err))
	}
	if overrideExp {
		claims.IssuedAt = jwt.NewNumericDate(time.Now())
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(100 * time.Hour))
	}
	fmt.Printf("\nAppend the following token into the 'Auth' header:\n\n%s\n\n",
		proxy.CreateRSASignedToken(claims, privateKey))

}
