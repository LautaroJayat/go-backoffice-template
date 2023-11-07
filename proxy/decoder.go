package proxy

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func DecodeToken(tokenString string, key *rsa.PublicKey) (*JWTStructure, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTStructure{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return key, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("there was an error while parsing claims and decoding token. error=%q", err)
	}

	claims, ok := token.Claims.(*JWTStructure)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("%s", "invalid token")
	}

	if claims.Role < 1 || claims.Email == "" || claims.Name == "" {
		return nil, fmt.Errorf("%s", "invalid token")
	}

	return claims, nil
}
