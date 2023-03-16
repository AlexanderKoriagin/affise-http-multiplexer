package main

import (
	"flag"
	"github.com/akrillis/affise-http-multiplexer/internal/check"
	"github.com/akrillis/affise-http-multiplexer/internal/handler"
	"github.com/akrillis/affise-http-multiplexer/internal/limit"
	"github.com/akrillis/affise-http-multiplexer/internal/server"
	"github.com/akrillis/affise-http-multiplexer/types"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	defaultPort              = "8080"
	defaultPath              = "/api/v1"
	defaultMaxIncome         = 100
	defaultMaxOutgoing       = 4
	defaultTimeoutStop       = 30 * time.Second
	defaultTimeoutPerRequest = 1 * time.Second
	defaultUrlQuantity       = 20
)

var (
	port              = flag.String("port", defaultPort, "HTTP listen port")
	path              = flag.String("path", defaultPath, "HTTP path")
	maxIncome         = flag.Uint64("maxIncome", defaultMaxIncome, "Quantity of incoming http requests working simultaneously")
	maxOutgoing       = flag.Int("maxOutgoing", defaultMaxOutgoing, "Quantity of outgoing http requests per each incoming request")
	timeoutStop       = flag.Duration("timeoutStop", defaultTimeoutStop, "Timeout for graceful shutdown, 20s for example")
	timeoutPerRequest = flag.Duration("timeoutPerRequest", defaultTimeoutPerRequest, "Timeout for each outgoing request, 1s for example")
	urlQuantity       = flag.Int("urlQuantity", defaultUrlQuantity, "Maximum quantity of urls in each incoming request")
)

func main() {
	flag.Parse()

	stop := make(chan os.Signal)
	swg := new(sync.WaitGroup)
	limiter := limit.Init(*maxIncome)
	checker := check.Init(*urlQuantity)
	h := handler.Init(*maxOutgoing, *timeoutPerRequest, checker)

	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	mux := http.NewServeMux()
	mux.HandleFunc(*path, limiter.RateHttp(h.CheckUrls))
	srv := server.New(
		types.ParamsServer{
			Mux:         mux,
			Port:        *port,
			TimeoutStop: *timeoutStop,
			Stop:        stop,
			Wait:        swg,
		},
	)

	swg.Add(1)
	srv.Start()

	log.Println("Server started")
	swg.Wait()
}
