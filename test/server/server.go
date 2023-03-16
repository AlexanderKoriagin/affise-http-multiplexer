package server

import (
	"github.com/akrillis/affise-http-multiplexer/internal/auxiliary"
	"net/http"
	"time"
)

const (
	PathOK   = "/ok"
	PathFail = "/fail"
	PathLong = "/long"
)

func ok(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Millisecond)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(http.StatusText(http.StatusOK)))
}

func ok2Seconds(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(http.StatusText(http.StatusOK)))
}

func fail(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Millisecond)
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

func Start() string {
	mux := http.NewServeMux()
	mux.HandleFunc(PathOK, ok)
	mux.HandleFunc(PathFail, fail)
	mux.HandleFunc(PathLong, ok2Seconds)

	addr := ":" + auxiliary.PortRandom()
	s := &http.Server{Handler: mux, Addr: addr}

	go func() {
		_ = s.ListenAndServe()
	}()

	return "http://" + addr
}
