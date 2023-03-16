package types

const (
	ErrServerStoppedWithError = "http server was stopped with error: %v\n"
	ErrServerUnableToStart    = "server unable to start: %v\n"
	ErrRequestContextDone     = "request context cancelled or exceeded"
	ErrTooManyUrls            = "too many urls"
	ErrOneOfRequestsFailed    = "one of requests failed"
)
