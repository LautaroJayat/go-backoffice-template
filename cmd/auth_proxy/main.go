package main

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lautarojayat/backoffice/proxy"
)

var (
	UPSTREAM_ENV   = "UPSTREAM"
	PUBLIC_KEY_ENV = "PUBLIC_KEY"
	PORT_ENV       = "PORT"
)

func main() {
	upstream := os.Getenv(UPSTREAM_ENV)
	pkey := os.Getenv(PUBLIC_KEY_ENV)
	port := os.Getenv(PORT_ENV)
	if upstream == "" || pkey == "" || port == "" {
		log.Fatal("missing config envs, cant continue")
	}

	parsedPKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pkey))
	if err != nil {
		log.Fatalf("could not parse public key. error=%q", err)
	}

	p, err := proxy.NewProxy(upstream)
	if err != nil {
		log.Fatalf("could not start proxy due to bad upstream url provided. error=%q", err)
	}
	http.HandleFunc("/", proxy.AuthChecker(p, parsedPKey))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
