package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/lautarojayat/backoffice/api/http/response"
	"github.com/lautarojayat/backoffice/roles"
)

func sendForbiddenError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	res, err := response.ToErrorResponse(response.BadRoleProvidedMsg, http.StatusForbidden)
	if err != nil {
		return
	}
	w.Write(res)
}

func CheckRole(l *log.Logger, roles []roles.Role, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		permString := r.Header.Get("X-Decoded-Perms")
		perms, err := strconv.ParseUint(permString, 2, 8)
		if err != nil {
			l.Printf("couldn't parse role. error=%q\n", err)
			sendForbiddenError(w)
			return
		}

		for _, role := range roles {
			if perms&uint64(role) == 0 {
				sendForbiddenError(w)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
