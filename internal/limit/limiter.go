package limit

import "net/http"

type Limiter struct {
	Features Features
}

func Init(tokens uint64) *Limiter {
	return &Limiter{Features: &service{max: tokens}}
}

func (l *Limiter) RateHttp(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !l.Features.Acquire() {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		h.ServeHTTP(w, r)
		l.Features.Release()
	}
}
