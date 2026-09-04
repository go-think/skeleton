package middleware

import (
	"github.com/go-think/flow"
)

// TrimStrings middleware trims whitespace from input
type TrimStrings struct{}

// NewTrimStrings creates TrimStrings middleware instance
func NewTrimStrings() flow.Handler {
	return &TrimStrings{}
}

// Process processes the request
func (m *TrimStrings) Process(req *flow.Request, next flow.Closure) interface{} {
	return next(req)
}
