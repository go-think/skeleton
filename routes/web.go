package routes

import (
	"net/http"

	"app/app/http/controllers"

	"github.com/go-think/flow"
	"github.com/go-think/think/filesystem"
)

func MapWebRoutes(r flow.Router) {
	welcome := controllers.NewWelcomeController()
	testCtrl := controllers.NewTestController()

	// 1. Static file service with directory listing protection
	staticHandler := http.StripPrefix("/static", http.FileServer(filesystem.NewFileFileSystem("./public", false)))
	r.Get("/static/*", staticHandler)

	// 2. Web Group with session
	r.Middleware("web").Group(func(web flow.Router) {
		web.Get("/", welcome.Index)
		web.Get("/status", welcome.Status)
		web.Get("/session/flow", testCtrl.SessionFlowAction)
	})

	// 3. View template rendering
	r.Get("/view/profile", testCtrl.RenderViewAction)

	// 4. Stream response and download
	r.Get("/stream", testCtrl.StreamAction)

	// 5. Exception handling test routes
	r.Get("/panic/http", testCtrl.PanicHttpExceptionAction)
	r.Get("/panic/server", testCtrl.PanicServerErrorAction)

	// 6. Redirect & NoContent routes
	r.Get("/redirect", func(req *flow.Request) *flow.Response {
		return flow.Redirect("/target", 301)
	})
	r.Get("/target", func(req *flow.Request) *flow.Response {
		return flow.Html("<h1>Target Reached</h1>")
	})
	r.Get("/nocontent", func(req *flow.Request) *flow.Response {
		return flow.NoContent()
	})

	// 7. Signed URL protected route
	r.Get("/signed/download", func(req *flow.Request) *flow.Response {
		return flow.Json(map[string]interface{}{"message": "secure content verified by signature"})
	}).Middleware(flow.NewValidateSignatureMiddleware(r)).Name("signed.download")

	// 8. Named route with parameters
	r.Get("/users/{id}", func(req *flow.Request) *flow.Response {
		id, _ := req.RouteParam("id")
		return flow.Json(map[string]interface{}{"user_id": id})
	}).Name("user.show")

	// 9. Route parameter regex constraints
	r.Get("/items/{id}/{code}/{status}", func(req *flow.Request) *flow.Response {
		id, _ := req.RouteParam("id")
		code, _ := req.RouteParam("code")
		status, _ := req.RouteParam("status")
		return flow.Json(map[string]interface{}{
			"id":     id,
			"code":   code,
			"status": status,
		})
	}).WhereNumber("id").WhereAlpha("code").WhereIn("status", []string{"active", "pending"})

	// 10. Fallback route (404)
	r.Fallback(func(req *flow.Request) *flow.Response {
		return flow.Json(map[string]interface{}{
			"error": "Custom Fallback: route not found",
			"path":  req.GetPath(),
		}).SetCode(404)
	})
}
