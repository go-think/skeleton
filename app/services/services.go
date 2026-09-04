package services

import (
	"fmt"
	"sync/atomic"
)

// UserService represents a singleton service
type UserService struct {
	Name string
}

func NewUserService() *UserService {
	return &UserService{Name: "CoreUserService"}
}

func (s *UserService) GetUser(id int) map[string]interface{} {
	return map[string]interface{}{
		"id":    id,
		"name":  fmt.Sprintf("User_%d", id),
		"email": fmt.Sprintf("user_%d@example.com", id),
	}
}

// OrderGenerator represents a transient binding (new instance every Make)
type OrderGenerator struct {
	ID int64
}

var orderSeq int64

func NewOrderGenerator() *OrderGenerator {
	return &OrderGenerator{
		ID: atomic.AddInt64(&orderSeq, 1),
	}
}

// RequestContext represents a scoped binding (per-request shared, flushed on terminate)
type RequestContext struct {
	TraceID string
}

func NewRequestContext() *RequestContext {
	return &RequestContext{
		TraceID: "trace-init",
	}
}

// PaymentService represents services for tagging
type PaymentService interface {
	Driver() string
}

type WechatPay struct{}

func (w *WechatPay) Driver() string { return "wechat" }

type AliPay struct{}

func (a *AliPay) Driver() string { return "alipay" }
