package server

import (
	"context"
	"github.com/akrillis/affise-http-multiplexer/types"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Features interface {
	Run() error
}

type server struct {
	mux     *http.ServeMux
	addr    string
	timeout time.Duration
	stop    chan os.Signal
	wait    *sync.WaitGroup
}

// Run - main function for route process which contains all rules for each handled URL
func (s *server) Run() error {
	srv := &http.Server{Handler: s.mux, Addr: s.addr}

	go func(hs *http.Server, cs chan os.Signal, wg *sync.WaitGroup) {
		<-cs
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
		defer cancel()

		if err := hs.Shutdown(ctx); err != nil {
			log.Printf(types.ErrServerStoppedWithError, err)
			return
		}
		log.Println(types.MsgServerStoppedGracefully)
	}(srv, s.stop, s.wait)

	return srv.ListenAndServe()
}
