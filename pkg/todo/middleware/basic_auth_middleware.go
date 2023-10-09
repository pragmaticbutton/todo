package middleware

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"todo/pkg/todo/dba"
)

type BasicAuthMiddleware struct {
	ErrorEncode func(http.ResponseWriter, error)
	Da          *dba.DatabaseAccess
}

func (bam BasicAuthMiddleware) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username, password, ok := r.BasicAuth()
		if !ok {
			bam.ErrorEncode(w, errors.New("basic auth is not correct"))
			return
		}

		u, err := bam.Da.SearchUser(nil, &username)
		if err != nil {
			bam.ErrorEncode(w, err)
			return
		}

		if len(u) == 0 {
			bam.ErrorEncode(w, errors.New("user not found"))
			return
		}

		s := sha256.New()
		s.Write([]byte(password))
		encr := fmt.Sprintf("%x", s.Sum(nil))

		if u[0].Password != encr {
			bam.ErrorEncode(w, errors.New("incorrect password"))
			return
		}

		h.ServeHTTP(w, r)
	})
}
