package middleware

import (
	"sync/atomic"

	"app/app/events"

	"github.com/go-think/flow"
)

var GlobalTerminatedCount int64

type AuditLogMiddleware struct{}

func NewAuditLogMiddleware() flow.Handler {
	return &AuditLogMiddleware{}
}

func (m *AuditLogMiddleware) Process(req *flow.Request, next flow.Closure) interface{} {
	return next(req)
}

func (m *AuditLogMiddleware) Terminate(req *flow.Request, res interface{}) {
	atomic.AddInt64(&GlobalTerminatedCount, 1)
	events.Tracker.Record("audit_terminated")
}
