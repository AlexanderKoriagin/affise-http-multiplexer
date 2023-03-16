package service

import "net/http"

type Limiter interface {
	RateHttp(h http.HandlerFunc) http.HandlerFunc
}
