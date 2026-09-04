package middleware

import (
	"fmt"
	"time"

	"github.com/go-think/flow"
)

type TimingMiddleware struct{}

func NewTimingMiddleware() flow.Handler {
	return &TimingMiddleware{}
}

func (m *TimingMiddleware) Process(req *flow.Request, next flow.Closure) interface{} {
	start := time.Now()
	req.Set("start_time", start)

	result := next(req)

	duration := time.Since(start)
	if res, ok := result.(*flow.Response); ok {
		res.Header.Set("X-Response-Time", fmt.Sprintf("%dms", duration.Milliseconds()))
	}
	return result
}
