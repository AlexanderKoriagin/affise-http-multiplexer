package handler

import (
	"context"
	"encoding/json"
	"github.com/akrillis/affise-http-multiplexer/service"
	"github.com/akrillis/affise-http-multiplexer/types"
	"log"
	"net/http"
	"sync"
	"time"
)

type Features interface {
	BaseRequest(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	parallel int           // simultaneous requests
	timeout  time.Duration // timeout for each request
	check    service.Checker
}

// BaseRequest
//  1. check if the request method is valid
//  2. check if the request body is valid
//  3. run the worker to process the request in parallel
//  4. return the result
func (h *handler) BaseRequest(w http.ResponseWriter, r *http.Request) {
	var (
		success = make(chan struct{}, 1)
		fail    = make(chan struct{}, 1)
		stop    = make(chan struct{}, 1)
		answer  types.AnswerBase
	)

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	var br types.RequestBase
	if err := json.NewDecoder(r.Body).Decode(&br); err != nil {
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	if !h.check.UrlsQty(br.Urls) {
		http.Error(w, types.ErrTooManyUrls, http.StatusBadRequest)
		return
	}

	go h.worker(r.Context(), br.Urls, stop, success, fail)

	select {
	case <-r.Context().Done():
		stop <- struct{}{}
		log.Println(types.ErrRequestContextDone)
		return
	case <-success:
		answer = types.AnswerBase{Result: types.MsgBaseRequestDone}
	case <-fail:
		stop <- struct{}{}
		http.Error(w, types.ErrOneOfRequestsFailed, http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(answer)
}

func (h *handler) worker(ctx context.Context, urls []string, stop, success, fail chan struct{}) {
	failed := make(chan struct{}, 1)
	sem := make(chan struct{}, h.parallel)
	wg := new(sync.WaitGroup)

	for _, url := range urls {
		select {
		case <-stop:
			return
		case <-failed:
			fail <- struct{}{}
			return
		default:
			sem <- struct{}{}

			ct, cancel := context.WithTimeout(ctx, h.timeout)
			defer cancel()

			wg.Add(1)
			go request(ct, url, failed, sem, wg)
		}
	}

	wg.Wait()

	select {
	case <-failed:
		fail <- struct{}{}
	default:
		success <- struct{}{}
	}
}

func request(ctx context.Context, url string, undone, sem chan struct{}, wg *sync.WaitGroup) {
	defer func() { <-sem }()
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		undone <- struct{}{}
		return
	}

	if resp, err := http.DefaultClient.Do(req); err != nil || resp.StatusCode != http.StatusOK {
		undone <- struct{}{}
	}
}
