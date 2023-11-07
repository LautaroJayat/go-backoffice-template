package proxy

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/lautarojayat/backoffice/roles"
)

type JWTStructure struct {
	Role  roles.Role `json:"role"`
	Name  string     `json:"name"`
	Email string     `json:"email"`
	jwt.RegisteredClaims
}
