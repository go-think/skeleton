package middleware

import (
	"net/http"

	"github.com/go-think/flow"
)

type RoleMiddleware struct {
	allowedRoles []string
}

func NewRoleMiddleware(roles ...string) flow.Handler {
	return &RoleMiddleware{allowedRoles: roles}
}

func (m *RoleMiddleware) Process(req *flow.Request, next flow.Closure) interface{} {
	userRole, _ := req.Query("role")
	if userRole == "" {
		userRole = req.Header("X-User-Role")
	}

	allowed := false
	for _, r := range m.allowedRoles {
		if r == userRole {
			allowed = true
			break
		}
	}

	if !allowed {
		return flow.Json(map[string]interface{}{
			"error": "Forbidden: insufficient role",
			"code":  http.StatusForbidden,
		}).SetCode(http.StatusForbidden)
	}

	return next(req)
}
