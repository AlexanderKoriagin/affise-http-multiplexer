package types

import (
	"net/http"
	"os"
	"sync"
	"time"
)

type ParamsServer struct {
	Mux         *http.ServeMux
	Port        string
	TimeoutStop time.Duration
	Stop        chan os.Signal
	Wait        *sync.WaitGroup
}
