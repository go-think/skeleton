package middleware

import (
	"net/http"

	"github.com/go-think/flow"
)

// Authenticate middleware validates authorization
type Authenticate struct{}

// NewAuthenticate creates an Authenticate middleware instance
func NewAuthenticate() flow.Handler {
	return &Authenticate{}
}

// Process handles request authentication inspection
func (m *Authenticate) Process(req *flow.Request, next flow.Closure) interface{} {
	token := req.Header("Authorization")
	if token == "" {
		token, _ = req.Query("token")
	}

	if token == "" {
		return flow.Json(map[string]interface{}{
			"error": "Unauthenticated.",
			"code":  http.StatusUnauthorized,
		}).SetCode(http.StatusUnauthorized)
	}

	return next(req)
}
