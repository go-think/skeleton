package controllers

import (
	"github.com/go-think/flow"
)

// BaseController provides common helper methods for all controllers
type BaseController struct{}

// Json returns a JSON response
func (c *BaseController) Json(data interface{}) *flow.Response {
	return flow.Json(data)
}

// Html returns an HTML response
func (c *BaseController) Html(content string) *flow.Response {
	return flow.Html(content)
}
