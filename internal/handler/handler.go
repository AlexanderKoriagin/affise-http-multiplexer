package handler

import (
	"github.com/akrillis/affise-http-multiplexer/service"
	"net/http"
	"time"
)

type Handler struct {
	Features Features
}

func Init(parallel int, timeout time.Duration, checker service.Checker) *Handler {
	return &Handler{
		Features: &handler{
			parallel: parallel,
			timeout:  timeout,
			check:    checker,
		},
	}
}

func (h *Handler) CheckUrls(w http.ResponseWriter, r *http.Request) {
	h.Features.BaseRequest(w, r)
}
