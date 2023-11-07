package proxy

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/lautarojayat/backoffice/roles"
)

func NewProxy(host string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return proxy, nil
}

func AuthChecker(proxy *httputil.ReverseProxy, publicKey *rsa.PublicKey) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Auth")
		if token == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		jwt, err := DecodeToken(token, publicKey)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(err.Error()))
			return
		}
		r.Header.Set(roles.DecodedPermsHeader, strconv.FormatUint(uint64(jwt.Role), 10))
		tokenString, _ := json.Marshal(jwt)
		r.Header.Set("Auth", string(tokenString))
		proxy.ServeHTTP(w, r)
	}
}
