package bootstrap

import (
	"fmt"
	"strings"

	"app/app/events"
	appMiddleware "app/app/http/middleware"
	"app/routes"

	"github.com/go-think/flow"
	"github.com/go-think/think"
)

// BootApp configures and builds the Think application
func BootApp() *think.Application {
	return think.Configure().
		WithRouting(think.Routing{
			Web:       routes.MapWebRoutes,
			Api:       routes.MapApiRoutes,
			ApiPrefix: "api",
			Health:    "/up",
		}).
		WithMiddleware(func(m *think.MiddlewareConfig) {
			// Global middlewares
			m.TrimStrings()
			m.Use(appMiddleware.NewTimingMiddleware())
			m.Use(appMiddleware.NewAuditLogMiddleware())
			m.Cors(flow.DefaultCorsConfig())

			// Middleware groups
			m.Group("web", flow.NewSessionMiddleware(nil))

			// Route-specific middleware aliases
			m.Alias("auth", appMiddleware.NewAuthenticate())
			m.Alias("role", func(roles ...string) flow.Handler {
				return appMiddleware.NewRoleMiddleware(roles...)
			})
		}).
		WithExceptions(func(e *think.ExceptionsConfig) {
			// Custom Report hook
			e.Report(func(err interface{}) bool {
				events.Tracker.Record(fmt.Sprintf("reported_err:%v", err))
				return false // continue default log report
			})

			// Custom Render hook
			e.Render(func(err interface{}) interface{} {
				errStr := fmt.Sprintf("%v", err)
				if strings.Contains(errStr, "database connection breakdown") {
					return flow.NewResponse().
						SetCode(500).
						SetContentType("application/json").
						SetContent(`{"error":"Database Service Unavailable","custom_rendered":true}`)
				}
				return nil // fallback to default render
			})
		}).
		WithProviders(AppProviders()...).
		Create()
}
