package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"app/app/services"

	"github.com/go-think/flow"
	"github.com/go-think/think/exception"
	"github.com/go-think/think/view"
)

type TestController struct {
	BaseController
}

func NewTestController() *TestController {
	return &TestController{}
}

// DependencyInjectionAction demonstrates auto-injecting Request and UserService from container
func (c *TestController) DependencyInjectionAction(req *flow.Request, userSvc *services.UserService) *flow.Response {
	svcName := "nil"
	if userSvc != nil {
		svcName = userSvc.Name
	}
	return c.Json(map[string]interface{}{
		"injected_service": svcName,
		"path":             req.GetPath(),
	})
}

// InputHelperAction demonstrates Request helper functions
func (c *TestController) InputHelperAction(req *flow.Request) *flow.Response {
	name, _ := req.Input("name", "Guest")
	age := req.Integer("age", 18)
	isAdmin := req.Boolean("is_admin", false)
	score := req.Float("score", 99.5)

	onlyMap := req.Only("name", "age")
	exceptMap := req.Except("token", "secret")
	segments := req.Segments()
	seg1 := req.Segment(1)

	req.Merge(map[string]string{"merged_key": "merged_val"})
	mergedVal, _ := req.Input("merged_key")

	return c.Json(map[string]interface{}{
		"name":        name,
		"age":         age,
		"is_admin":    isAdmin,
		"score":       score,
		"only":        onlyMap,
		"except":      exceptMap,
		"segments":    segments,
		"seg1":        seg1,
		"merged_val":  mergedVal,
		"has_name":    req.Has("name"),
		"exists_name": req.Exists("name"),
		"client_ip":   req.ClientIP(),
		"bearer":      req.BearerToken(),
		"url":         req.Url(),
		"full_url":    req.FullUrl(),
	})
}

// SessionFlowAction demonstrates Session Flash, Keep, Reflash, Token, Increment
func (c *TestController) SessionFlowAction(req *flow.Request) *flow.Response {
	sess := req.Session()
	if sess == nil {
		return c.Json(map[string]interface{}{"error": "session not found"})
	}

	action, _ := req.Query("action")
	switch action {
	case "set_flash":
		sess.Flash("status_msg", "Order Created Successfully")
		sess.Set("user_id", 1001)
		sess.Increment("counter", 1)
		sess.SetPreviousUrl("/orders")
		return c.Json(map[string]interface{}{
			"message": "flash set",
			"token":   sess.Token(),
		})
	case "get_flash":
		msg := sess.Get("status_msg")
		counter := sess.Get("counter")
		prev := sess.PreviousUrl()
		return c.Json(map[string]interface{}{
			"status_msg": msg,
			"counter":    counter,
			"previous":   prev,
		})
	case "reflash":
		sess.Reflash()
		return c.Json(map[string]interface{}{"message": "reflashed"})
	case "regenerate_token":
		newToken := sess.RegenerateToken()
		return c.Json(map[string]interface{}{"token": newToken})
	default:
		return c.Json(map[string]interface{}{"message": "idle"})
	}
}

// FileUploadAction demonstrates UploadedFile and Move
func (c *TestController) FileUploadAction(req *flow.Request) *flow.Response {
	file, err := req.File("avatar")
	if err != nil || file == nil {
		return c.Json(map[string]interface{}{"error": "no file uploaded"}).SetCode(http.StatusBadRequest)
	}

	saveDir := "storage/app/uploads"
	_ = os.MkdirAll(saveDir, 0755)
	ok, err := file.Move(saveDir, "uploaded_avatar.txt")
	if err == nil && ok {
		return c.Json(map[string]interface{}{
			"message": "file uploaded successfully",
			"saved":   filepath.Join(saveDir, "uploaded_avatar.txt"),
		})
	}
	return c.Json(map[string]interface{}{"error": fmt.Sprintf("%v", err)}).SetCode(http.StatusInternalServerError)
}

// StreamAction demonstrates Server-Sent Events (SSE) and Streamed downloads
func (c *TestController) StreamAction(req *flow.Request) *flow.Response {
	mode, _ := req.Query("mode")
	if mode == "download" {
		return flow.StreamDownload(func(w io.Writer) bool {
			_, _ = w.Write([]byte("id,name,role\n1,Alice,Admin\n2,Bob,User\n"))
			return false
		}, "export.csv")
	}

	count := 0
	return flow.StreamResponse(func(w io.Writer) bool {
		count++
		_, _ = fmt.Fprintf(w, "data: event-%d\n\n", count)
		return count < 3
	})
}

// PanicHttpExceptionAction tests HttpException (e.g. 403)
func (c *TestController) PanicHttpExceptionAction(req *flow.Request) *flow.Response {
	panic(exception.NewHttpException(http.StatusForbidden, "Access Denied By Security Rule"))
}

// PanicServerErrorAction tests unhandled native Panic
func (c *TestController) PanicServerErrorAction(req *flow.Request) *flow.Response {
	panic("unexpected database connection breakdown")
}

// RenderViewAction demonstrates template rendering via think.View
func (c *TestController) RenderViewAction(req *flow.Request) *flow.Response {
	html := view.Render("resources/views/profile.html", map[string]interface{}{
		"Name":  "Zhang San",
		"Email": "zhangsan@example.com",
	})
	return flow.Html(string(html))
}
