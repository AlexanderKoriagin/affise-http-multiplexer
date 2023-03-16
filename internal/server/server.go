package server

import (
	"github.com/akrillis/affise-http-multiplexer/types"
	"log"
	"net/http"
)

type Server struct {
	Features Features
}

func New(params types.ParamsServer) *Server {
	return &Server{
		Features: &server{
			mux:     params.Mux,
			addr:    ":" + params.Port,
			timeout: params.TimeoutStop,
			stop:    params.Stop,
			wait:    params.Wait,
		},
	}
}

func (s *Server) Start() {
	go func() {
		if err := s.Features.Run(); err != nil && err != http.ErrServerClosed {
			log.Printf(types.ErrServerUnableToStart, err)
		}
	}()
}
