package middleware

import (
	"net/http"

	"github.com/go-think/flow"
)

// Authenticate middleware validates session-based or token-based authorization
type Authenticate struct{}

// NewAuthenticate creates an Authenticate middleware instance
func NewAuthenticate() flow.Handler {
	return &Authenticate{}
}

// Process handles request authentication inspection
func (m *Authenticate) Process(req *flow.Request, next flow.Closure) interface{} {
	// 1. Session-based authentication (Web requests)
	if s := req.Session(); s != nil && s.Has("user_id") {
		return next(req)
	}

	// 2. Token-based authentication (API requests)
	token := req.Header("Authorization")
	if token == "" {
		token, _ = req.Query("token")
	}

	if token != "" {
		return next(req)
	}

	// 3. Unauthenticated response
	if req.ExpectsJson() {
		return flow.Json(map[string]interface{}{
			"error": "Unauthenticated.",
			"code":  http.StatusUnauthorized,
		}).SetCode(http.StatusUnauthorized)
	}

	return flow.Redirect("/login")
}
