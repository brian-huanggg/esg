package middleware

import (
	"context"
	"crypto/rand"
	"net/http"
	"regexp"
)

// HeaderRequestID is the header used to read and return the request id.
const HeaderRequestID = "X-Request-Id"

// ctxKey is unexported so no other package can read or overwrite our value.
type ctxKey struct{}

// Client-supplied ids end up in logs and response headers, so only accept a
// short, safe charset; anything else is replaced.
var validID = regexp.MustCompile(`^[A-Za-z0-9-]{1,64}$`)

// RequestID reuses a valid incoming X-Request-Id or generates a new one,
// stores it in the request context, and echoes it in the response header.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(HeaderRequestID)
		if !validID.MatchString(id) {
			id = rand.Text()
		}
		w.Header().Set(HeaderRequestID, id)
		ctx := context.WithValue(r.Context(), ctxKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFrom returns the request id, or "" if the middleware did not run.
func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}
