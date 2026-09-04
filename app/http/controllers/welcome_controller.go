package controllers

import (
	"os"
	"runtime"
	"time"

	"app/config"

	"github.com/go-think/flow"
	"github.com/go-think/think/facades"
)

// WelcomeController handles default landing and status endpoints
type WelcomeController struct {
	BaseController
}

// NewWelcomeController creates WelcomeController instance
func NewWelcomeController() *WelcomeController {
	return &WelcomeController{}
}

// Index renders the welcome landing page
func (c *WelcomeController) Index(req *flow.Request) *flow.Response {
	viewContent, err := os.ReadFile("resources/views/welcome.html")
	if err == nil {
		return flow.Html(string(viewContent))
	}
	return flow.Html("<h1>Welcome to Think Framework!</h1>")
}

// Status returns runtime server status in JSON
func (c *WelcomeController) Status(req *flow.Request) *flow.Response {
	session := req.Session()
	var visitCount int = 1
	if session != nil {
		if val := session.Get("visit_count"); val != nil {
			if count, ok := val.(int); ok {
				visitCount = count + 1
			}
		}
		session.Set("visit_count", visitCount)
	}

	return c.Json(map[string]interface{}{
		"app":         config.App.Name,
		"version":     "1.0.0",
		"status":      "running",
		"go_version":  runtime.Version(),
		"server_time": time.Now().Format(time.RFC3339),
		"visit_count": visitCount,
		"env":         config.App.Env,
		"debug":       config.App.Debug,
		"is_local":    facades.App.IsLocal(),
		"is_prod":     facades.App.IsProduction(),
		"db_default":  config.Database.Default,
		"db_host":     config.Database.Connections.MySQL.Host,
	})
}
