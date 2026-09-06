package routes

import (
	"app/app/http/controllers"

	"github.com/go-think/flow"
)

func MapApiRoutes(r flow.Router) {
	testCtrl := controllers.NewTestController()

	r.Group(flow.GroupAttributes{Middleware: []interface{}{"api"}}, func(api flow.Router) {
		api.Get("/ping", func(req *flow.Request) *flow.Response {
			return flow.Json(map[string]string{"message": "pong"})
		})

		// Protected by 'auth' alias
		api.Get("/user", func(req *flow.Request) *flow.Response {
			return flow.Json(map[string]interface{}{
				"id":    1,
				"name":  "Think Developer",
				"email": "dev@go-think.com",
			})
		}).Middleware("auth")

		// Protected by parameterized 'role:admin,manager' alias
		api.Get("/admin", func(req *flow.Request) *flow.Response {
			return flow.Json(map[string]interface{}{
				"message": "welcome admin",
			})
		}).Middleware("role:admin,manager")

		// Dependency Injection testing
		api.Get("/di", testCtrl.DependencyInjectionAction)

		// Input helper testing
		api.Get("/input", testCtrl.InputHelperAction)
		api.Post("/input", testCtrl.InputHelperAction)

		// File upload testing
		api.Post("/upload", testCtrl.FileUploadAction)

		// Verbs testing
		api.Post("/resource", func(req *flow.Request) *flow.Response {
			return flow.Json(map[string]string{"action": "create"}).SetCode(201)
		})
		api.Put("/resource/{id}", func(req *flow.Request) *flow.Response {
			id, _ := req.RouteParam("id")
			return flow.Json(map[string]string{"action": "update", "id": id})
		})
		api.Patch("/resource/{id}", func(req *flow.Request) *flow.Response {
			id, _ := req.RouteParam("id")
			return flow.Json(map[string]string{"action": "patch", "id": id})
		})
		api.Delete("/resource/{id}", func(req *flow.Request) *flow.Response {
			id, _ := req.RouteParam("id")
			return flow.Json(map[string]string{"action": "delete", "id": id})
		})
		api.Any("/any", func(req *flow.Request) *flow.Response {
			return flow.Json(map[string]string{"method": req.Method()})
		})
	})
}
