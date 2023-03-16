package server

import (
	"fmt"
	"github.com/akrillis/affise-http-multiplexer/internal/auxiliary"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

func TestServer_Run(t *testing.T) {
	mux := http.NewServeMux()
	stop := make(chan os.Signal)
	wg := new(sync.WaitGroup)

	s := &server{
		mux:     mux,
		timeout: 10 * time.Second,
		stop:    stop,
		wait:    wg,
	}

	// wrong port
	s.addr = ":99999"
	assert.NotNil(t, s.Run())

	// run and stop successfully
	s.addr = ":" + auxiliary.PortRandom()
	c := make(chan error, 1)

	wg.Add(1)
	go func(c chan error) {
		if err := s.Run(); err != nil {
			c <- fmt.Errorf("unable to start http server: %v", err)
		}
	}(c)

	ts := time.After(300 * time.Millisecond)

	select {
	case err := <-c:
		log.Fatalln(err)
	case <-ts:
		log.Println("http server started")
	}

	stop <- os.Interrupt
	time.Sleep(100 * time.Millisecond)
}
