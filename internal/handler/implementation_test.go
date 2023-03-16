package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/akrillis/affise-http-multiplexer/internal/check"
	testServer "github.com/akrillis/affise-http-multiplexer/test/server"
	"github.com/akrillis/affise-http-multiplexer/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHandler_BaseRequest(t *testing.T) {
	addr := testServer.Start()

	h := &handler{
		parallel: 2,
		timeout:  200 * time.Millisecond,
		check:    check.Init(64),
	}

	correctUrls, incorrectUrls := []string(nil), []string(nil)
	correctUrlsButOneTooLong, tooManyUrls := []string(nil), []string(nil)

	for i := 0; i < 32; i++ {
		correctUrls = append(correctUrls, addr+testServer.PathOK)
		correctUrlsButOneTooLong = append(correctUrls, addr+testServer.PathOK)
		incorrectUrls = append(incorrectUrls, addr+testServer.PathOK)
	}
	incorrectUrls[16] = addr + testServer.PathFail
	correctUrlsButOneTooLong[16] = addr + testServer.PathLong

	for i := 0; i < 128; i++ {
		tooManyUrls = append(tooManyUrls, addr+testServer.PathOK)
	}

	//
	// incorrect http method
	//
	r := httptest.NewRequest(http.MethodPut, "/", nil)
	w := httptest.NewRecorder()
	h.BaseRequest(w, r)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	//
	// incorrect body
	//
	body, err := json.Marshal("incorrect body")
	assert.Nil(t, err)
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	h.BaseRequest(w, r)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	//
	// incorrect urls qty
	//
	body, err = json.Marshal(types.RequestBase{Urls: tooManyUrls})
	assert.Nil(t, err)
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	h.BaseRequest(w, r)
	assert.Equal(t, types.ErrTooManyUrls+"\n", w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)

	body, err = json.Marshal(types.RequestBase{})
	assert.Nil(t, err)
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	h.BaseRequest(w, r)
	assert.Equal(t, types.ErrTooManyUrls+"\n", w.Body.String())
	assert.Equal(t, http.StatusBadRequest, w.Code)

	//
	// request context timeout exceeded
	//
	body, err = json.Marshal(types.RequestBase{Urls: correctUrls})
	assert.Nil(t, err)
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	h.BaseRequest(w, r.WithContext(ctx))

	//
	// failed execution
	//
	body, err = json.Marshal(types.RequestBase{Urls: incorrectUrls})
	assert.Nil(t, err)
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h.BaseRequest(w, r.WithContext(ctx))
	assert.Equal(t, types.ErrOneOfRequestsFailed+"\n", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	//
	// failed execution because of too long outgoing request
	//
	body, err = json.Marshal(types.RequestBase{Urls: correctUrlsButOneTooLong})
	assert.Nil(t, err)
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h.BaseRequest(w, r.WithContext(ctx))
	assert.Equal(t, types.ErrOneOfRequestsFailed+"\n", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	//
	// success execution
	//
	body, err = json.Marshal(types.RequestBase{Urls: correctUrls})
	assert.Nil(t, err)
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h.BaseRequest(w, r.WithContext(ctx))
	assert.Equal(t, "{\"result\":\"done\"}\n", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_worker(t *testing.T) {
	addr := testServer.Start()

	success := make(chan struct{}, 1)
	fail := make(chan struct{}, 1)
	stop := make(chan struct{}, 1)

	h := &handler{
		parallel: 2,
		timeout:  200 * time.Millisecond,
	}

	correctUrls, incorrectUrls := []string(nil), []string(nil)
	for i := 0; i < 16; i++ {
		correctUrls = append(correctUrls, addr+testServer.PathOK)
		incorrectUrls = append(incorrectUrls, addr+testServer.PathOK)
	}
	incorrectUrls[8] = addr + testServer.PathFail

	//
	// context timeout exceeded
	//
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	go h.worker(ctx, correctUrls, stop, success, fail)
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 0, len(success))
	assert.Equal(t, 1, len(fail))
	<-fail

	//
	// failed execution
	//
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go h.worker(ctx, incorrectUrls, stop, success, fail)
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 0, len(success))
	assert.Equal(t, 1, len(fail))
	<-fail

	//
	// success execution
	//
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go h.worker(ctx, correctUrls, stop, success, fail)
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 1, len(success))
	assert.Equal(t, 0, len(fail))
}

func TestHandler_request(t *testing.T) {
	addr := testServer.Start()

	undone := make(chan struct{}, 1)
	sem := make(chan struct{}, 1)
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	//
	// test request with timeout
	//
	sem <- struct{}{}
	wg.Add(1)
	go request(ctx, addr+testServer.PathOK, undone, sem, wg)
	wg.Wait()
	assert.Equal(t, 0, len(sem))
	assert.Equal(t, 1, len(undone))
	<-undone

	//
	// test request with success
	//
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sem <- struct{}{}
	wg.Add(1)
	go request(ctx, addr+testServer.PathOK, undone, sem, wg)
	wg.Wait()
	assert.Equal(t, 0, len(sem))
	assert.Equal(t, 0, len(undone))

	//
	// test request with fail
	//
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sem <- struct{}{}
	wg.Add(1)
	go request(ctx, addr+testServer.PathFail, undone, sem, wg)
	wg.Wait()
	assert.Equal(t, 0, len(sem))
	assert.Equal(t, 1, len(undone))
}
