package limit

import "sync/atomic"

type Features interface {
	Acquire() bool
	Release()
}

type service struct {
	max    uint64 // maximum number of tokens
	tokens uint64 // current number of tokens
}

func (s *service) Acquire() bool {
	if atomic.LoadUint64(&s.tokens) < s.max {
		atomic.AddUint64(&s.tokens, 1)
		return true
	}

	return false
}

func (s *service) Release() {
	if s.tokens > 0 {
		atomic.AddUint64(&s.tokens, ^uint64(0))
	}
}
