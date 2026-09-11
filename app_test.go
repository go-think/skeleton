package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"app/app/events"
	"app/app/services"
	"app/bootstrap"
	"app/config"

	"github.com/go-think/flow"
	"github.com/go-think/think/contract"
	"github.com/go-think/think/facades"
	"github.com/go-think/think/helper"
	"github.com/stretchr/testify/assert"
)

// 1. Test Application Lifecycle and IoC Container
func TestApplicationAndContainer(t *testing.T) {
	app := bootstrap.BootApp()
	assert.NotNil(t, app)

	// 1.1 Singleton binding test
	userSvc1 := app.Make[*services.UserService]("UserService")
	userSvc2 := app.Make[*services.UserService]("UserService")
	assert.Same(t, userSvc1, userSvc2, "Singletons must return exact same instance pointer")
	assert.Equal(t, "CoreUserService", userSvc1.Name)

	// 1.2 Transient (Bind) test
	orderGen1 := app.Make[*services.OrderGenerator]("OrderGenerator")
	orderGen2 := app.Make[*services.OrderGenerator]("OrderGenerator")
	assert.NotSame(t, orderGen1, orderGen2, "Transient bindings must return newly instantiated objects")
	assert.NotEqual(t, orderGen1.ID, orderGen2.ID, "Sequential IDs must be distinct")

	// 1.3 Scoped binding test
	scoped1 := app.Make[*services.RequestContext]("RequestContext")
	scoped2 := app.Make[*services.RequestContext]("RequestContext")
	assert.Same(t, scoped1, scoped2, "Scoped bindings must be shared before flush")

	app.FlushScoped()
	scoped3 := app.Make[*services.RequestContext]("RequestContext")
	assert.NotSame(t, scoped1, scoped3, "Scoped bindings must be re-created after flush")

	// 1.4 Tagged bindings test
	taggedObjs := app.Tagged[services.PaymentService]("payments")
	assert.Len(t, taggedObjs, 2, "Tagged payments should resolve all 2 driver instances")
	var drivers []string
	for _, obj := range taggedObjs {
		drivers = append(drivers, obj.Driver())
	}
	assert.Contains(t, drivers, "wechat")
	assert.Contains(t, drivers, "alipay")

	// 1.5 Reflection Call with Auto-wiring test
	callResults := app.Call(func(svc *services.UserService) string {
		return "Hello " + svc.Name
	})
	assert.Len(t, callResults, 1)
	assert.Equal(t, "Hello CoreUserService", callResults[0])
}

// 2. Test Config Repository, Environment and Helpers
func TestConfigAndEnvironment(t *testing.T) {
	app := bootstrap.BootApp()
	_ = app

	// 2.1 Facades config access
	cfg := facades.Config()
	assert.Equal(t, "Think", cfg.GetString("app.name"))
	assert.Equal(t, "8080", cfg.GetString("app.port"))
	assert.Equal(t, 3306, cfg.GetInt("database.connections.mysql.port"))
	assert.True(t, cfg.GetBool("app.debug"))
	assert.NotEmpty(t, cfg.GetString("logging.default"))
	assert.NotNil(t, config.Logging)
	assert.NotEmpty(t, config.Logging.Default)

	// 2.2 Array config manipulation: Push and Prepend
	cfg.Push("app.providers_list", "ProviderA")
	cfg.Push("app.providers_list", "ProviderB")
	cfg.Prepend("app.providers_list", "ProviderHead")
	providers := cfg.Get("app.providers_list").([]interface{})
	assert.Equal(t, "ProviderHead", providers[0])
	assert.Equal(t, "ProviderA", providers[1])
	assert.Equal(t, "ProviderB", providers[2])

	// 2.3 Helper functions
	assert.Equal(t, "Think", helper.Config("app.name"))
	assert.NotEmpty(t, helper.WorkPath())
	assert.NotEmpty(t, helper.ConfigPath())
	assert.NotEmpty(t, helper.AppPath())

	// 2.4 Environment assertions
	assert.True(t, app.IsLocal())
	assert.False(t, facades.App.IsProduction())
	assert.Equal(t, "UTC", facades.App.GetTimezone())
}

// 3. Test HTTP Endpoints, Routing Verbs, Regex Constraints & Fallback
func TestHttpEndpointsAndRouting(t *testing.T) {
	app := bootstrap.BootApp()
	server := httptest.NewServer(app.BuildServer().Handler)
	defer server.Close()

	// 3.1 Health probe GET /up
	resp, err := http.Get(server.URL + "/up")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "UP", string(body))
	resp.Body.Close()

	// 3.2 HTTP Verbs testing (POST, PUT, PATCH, DELETE, ANY)
	// POST /api/resource
	resp, err = http.Post(server.URL+"/api/resource", "application/json", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "create")
	resp.Body.Close()

	// PUT /api/resource/42
	req, _ := http.NewRequest("PUT", server.URL+"/api/resource/42", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "update")
	assert.Contains(t, string(body), "42")
	resp.Body.Close()

	// PATCH /api/resource/42
	req, _ = http.NewRequest("PATCH", server.URL+"/api/resource/42", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "patch")
	resp.Body.Close()

	// DELETE /api/resource/42
	req, _ = http.NewRequest("DELETE", server.URL+"/api/resource/42", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "delete")
	resp.Body.Close()

	// ANY /api/any
	req, _ = http.NewRequest("GET", server.URL+"/api/any", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "GET")
	resp.Body.Close()

	// 3.3 Route Parameter Regex Constraints: WhereNumber, WhereAlpha, WhereIn
	// Match: id=123 (num), code=abc (alpha), status=active (in allowed list)
	resp, err = http.Get(server.URL + "/items/123/abc/active")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "123")
	assert.Contains(t, string(body), "abc")
	assert.Contains(t, string(body), "active")
	resp.Body.Close()

	// Mismatch 1: non-numeric ID -> triggers Fallback 404
	resp, err = http.Get(server.URL + "/items/invalid-num/abc/active")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Custom Fallback: route not found")
	resp.Body.Close()

	// Mismatch 2: non-alpha code -> triggers Fallback 404
	resp, err = http.Get(server.URL + "/items/123/12345/active")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	// Mismatch 3: invalid status not in [active, pending] -> triggers Fallback 404
	resp, err = http.Get(server.URL + "/items/123/abc/archived")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	// 3.4 Named Route & URL Generation
	r := app.Make[flow.Router]("router")
	generatedUrl := r.Url("user.show", map[string]string{"id": "999"})
	assert.Equal(t, "/users/999", generatedUrl)

	resp, err = http.Get(server.URL + generatedUrl)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "999")
	resp.Body.Close()

	// 3.5 Controller Dependency Injection via Route Dispatcher
	resp, err = http.Get(server.URL + "/api/di")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "CoreUserService", "UserService should be auto-injected into controller action")
	resp.Body.Close()

	// 3.6 Redirect & NoContent
	req, _ = http.NewRequest("GET", server.URL+"/redirect", nil)
	// Disable auto-redirect to verify 301 status code
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode)
	assert.Equal(t, "/target", resp.Header.Get("Location"))
	resp.Body.Close()

	resp, err = http.Get(server.URL + "/nocontent")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()
}

// 4. Test Signed URL & Security Protection
func TestSignedUrlAndSecurity(t *testing.T) {
	// URL signing requires a configured app.key; without it flow fails closed.
	t.Setenv("APP_KEY", "test-signature-key")

	app := bootstrap.BootApp()
	server := httptest.NewServer(app.BuildServer().Handler)
	defer server.Close()

	// Route is pre-registered in routes/web.go with Name("signed.download")

	// 4.1 Unsigned request -> 403 Forbidden
	resp, err := http.Get(server.URL + "/signed/download")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid signature")
	resp.Body.Close()

	// 4.2 Valid signature -> 200 OK
	urlGenerator := app.Make[*flow.UrlGenerator]()
	signedUrl := urlGenerator.TemporarySignedRoute("signed.download", 5*time.Minute, nil)
	assert.Contains(t, signedUrl, "signature=")
	assert.Contains(t, signedUrl, "expires=")

	resp, err = http.Get(server.URL + signedUrl)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "secure content verified by signature")
	resp.Body.Close()

	// 4.3 Tampered signature -> 403 Forbidden
	tamperedUrl := signedUrl + "&tampered=true"
	resp, err = http.Get(server.URL + tamperedUrl)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	resp.Body.Close()

	// 4.4 Expired signature -> 403 Forbidden
	expiredUrl := urlGenerator.TemporarySignedRoute("signed.download", -1*time.Minute, nil)
	resp, err = http.Get(server.URL + expiredUrl)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	resp.Body.Close()
}

// 5. Test Middlewares: Timing, CORS, Parameterized Role, and Terminable
func TestMiddlewareAndPipelines(t *testing.T) {
	app := bootstrap.BootApp()
	server := httptest.NewServer(app.BuildServer().Handler)
	defer server.Close()

	events.Tracker.Reset()

	// 5.1 Onion Model & Timing Middleware (X-Response-Time header)
	resp, err := http.Get(server.URL + "/api/ping")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("X-Response-Time"), "Response should include custom timing middleware header")
	resp.Body.Close()

	// 5.2 CORS Middleware handling preflight OPTIONS
	req, _ := http.NewRequest("OPTIONS", server.URL+"/api/ping", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Methods"), "POST")
	resp.Body.Close()

	// 5.3 Parameterized Middleware: role:admin,manager
	// Denied (no role provided)
	resp, err = http.Get(server.URL + "/api/admin")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	resp.Body.Close()

	// Allowed with admin role
	req, _ = http.NewRequest("GET", server.URL+"/api/admin?role=admin", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "welcome admin")
	resp.Body.Close()

	// Allowed with manager role
	req, _ = http.NewRequest("GET", server.URL+"/api/admin?role=manager", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Denied with regular user role
	req, _ = http.NewRequest("GET", server.URL+"/api/admin?role=user", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	resp.Body.Close()

	// 5.4 Terminable Middleware and Kernel Events (Executed asynchronously)
	// Give asynchronous Goroutine a brief moment to complete
	time.Sleep(50 * time.Millisecond)
	assert.True(t, events.Tracker.Has("audit_terminated"), "Terminable middleware should be executed asynchronously")
	assert.True(t, events.Tracker.Has("kernel.handled"), "kernel.handled event should be dispatched")
	assert.True(t, events.Tracker.Has("kernel.terminating"), "kernel.terminating event should be dispatched")
}

// 6. Test Flow Request Input Helpers, File Upload & Streaming Response
func TestFlowRequestAndInputHelpers(t *testing.T) {
	app := bootstrap.BootApp()
	server := httptest.NewServer(app.BuildServer().Handler)
	defer server.Close()

	// 6.1 Input parsing, Boolean/Integer/Float casting, TrimStrings
	// Pass leading/trailing spaces to test TrimStrings middleware
	reqUrl := server.URL + "/api/input?name=%20%20JohnDoe%20%20&age=25&is_admin=true&score=98.5&token=secret-token"
	resp, err := http.Get(reqUrl)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	assert.Contains(t, bodyStr, `"name":"JohnDoe"`, "TrimStrings middleware should have stripped whitespace")
	assert.Contains(t, bodyStr, `"age":25`)
	assert.Contains(t, bodyStr, `"is_admin":true`)
	assert.Contains(t, bodyStr, `"score":98.5`)
	assert.Contains(t, bodyStr, `"merged_val":"merged_val"`)
	assert.Contains(t, bodyStr, `"has_name":true`)
	assert.Contains(t, bodyStr, `"exists_name":true`)
	resp.Body.Close()

	// 6.2 Multipart File Upload and Move
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, _ := w.CreateFormFile("avatar", "test_avatar.txt")
	_, _ = part.Write([]byte("Hello avatar file content!"))
	_ = w.Close()

	req, _ := http.NewRequest("POST", server.URL+"/api/upload", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "file uploaded successfully")
	resp.Body.Close()

	// Verify saved file on disk
	savedContent, err := os.ReadFile("storage/app/uploads/uploaded_avatar.txt")
	assert.NoError(t, err)
	assert.Equal(t, "Hello avatar file content!", string(savedContent))

	// 6.3 Server-Sent Events (SSE) Streaming Response
	resp, err = http.Get(server.URL + "/stream")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/event-stream")
	body, _ = io.ReadAll(resp.Body)
	streamContent := string(body)
	assert.Contains(t, streamContent, "data: event-1")
	assert.Contains(t, streamContent, "data: event-2")
	assert.Contains(t, streamContent, "data: event-3")
	resp.Body.Close()

	// 6.4 Streamed File Download
	resp, err = http.Get(server.URL + "/stream?mode=download")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Disposition"), "export.csv")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Alice,Admin")
	resp.Body.Close()
}

// 7. Test Session Full Lifecycle: Flash Aging, Counter, Previous URL & CSRF Token
func TestSessionFullLifecycle(t *testing.T) {
	app := bootstrap.BootApp()
	server := httptest.NewServer(app.BuildServer().Handler)
	defer server.Close()

	// Use cookiejar to simulate browser session persistence
	var sessionCookie *http.Cookie

	// Step 1: Set Flash message and counter
	req, _ := http.NewRequest("GET", server.URL+"/session/flow?action=set_flash", nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	for _, c := range resp.Cookies() {
		if c.Name == config.Session.Cookie {
			sessionCookie = c
			break
		}
	}
	assert.NotNil(t, sessionCookie, "Session cookie must be set in response")
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "flash set")
	resp.Body.Close()

	// Step 2: Immediate next request -> Flash message should be available (moved to old)
	req, _ = http.NewRequest("GET", server.URL+"/session/flow?action=get_flash", nil)
	req.AddCookie(sessionCookie)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	body, _ = io.ReadAll(resp.Body)
	bodyStr := string(body)
	assert.Contains(t, bodyStr, "Order Created Successfully", "Flash message must be available on 1st follow-up request")
	assert.Contains(t, bodyStr, `"counter":1`)
	assert.Contains(t, bodyStr, `"/orders"`)
	resp.Body.Close()

	// Step 3: Second next request without reflash -> Flash message must be expired and absent!
	req, _ = http.NewRequest("GET", server.URL+"/session/flow?action=get_flash", nil)
	req.AddCookie(sessionCookie)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	body, _ = io.ReadAll(resp.Body)
	bodyStr = string(body)
	assert.NotContains(t, bodyStr, "Order Created Successfully", "Flash message must be expired on 2nd follow-up request")
	assert.Contains(t, bodyStr, `"status_msg":null`)
	resp.Body.Close()

	// Step 4: Regenerate CSRF Token
	req, _ = http.NewRequest("GET", server.URL+"/session/flow?action=regenerate_token", nil)
	req.AddCookie(sessionCookie)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "token")
	resp.Body.Close()
}

// 8. Test Exception Handling, Recover, and Custom Render Hooks
func TestExceptionHandlingAndRecovery(t *testing.T) {
	app := bootstrap.BootApp()
	server := httptest.NewServer(app.BuildServer().Handler)
	defer server.Close()

	events.Tracker.Reset()

	// 8.1 Panic with HttpException(403)
	resp, err := http.Get(server.URL + "/panic/http")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Access Denied By Security Rule")
	resp.Body.Close()

	// 8.2 Panic with native error intercepted by Custom Render Hook -> 500 JSON
	resp, err = http.Get(server.URL + "/panic/server")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Database Service Unavailable")
	assert.Contains(t, string(body), `"custom_rendered":true`)
	resp.Body.Close()

	// 8.3 Custom Report Hook execution assertion
	assert.True(t, events.Tracker.Has("reported_err:unexpected database connection breakdown"))
}

// 9. Test View Rendering via template.HTML
func TestViewRendering(t *testing.T) {
	app := bootstrap.BootApp()
	server := httptest.NewServer(app.BuildServer().Handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/view/profile")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	htmlStr := string(body)
	assert.Contains(t, htmlStr, "User: Zhang San")
	assert.Contains(t, htmlStr, "Email: zhangsan@example.com")
	resp.Body.Close()
}

// 10. Test Static File Protection (Direct FileServer Stream)
func TestStaticFileProtection(t *testing.T) {
	app := bootstrap.BootApp()
	server := httptest.NewServer(app.BuildServer().Handler)
	defer server.Close()

	// 10.1 Verify static file content is truly streamed to client
	resp, err := http.Get(server.URL + "/static/static.txt")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "This is static asset content for testing.", "Static file content must be served directly")
	resp.Body.Close()

	// 10.2 Verify directory listing protection (neuteredStatFile returns 404)
	resp, err = http.Get(server.URL + "/static/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Directory without index.html must return 404")
	resp.Body.Close()
}

// 11. Test Event Dispatcher (Listen, Dispatch, Until)
func TestEventDispatcher(t *testing.T) {
	app := bootstrap.BootApp()
	app.Boot() // Explicitly trigger service providers boot in non-server tests
	dispatcher := app.Make[contract.EventDispatcher]("contract.EventDispatcher")

	events.Tracker.Reset()

	// 11.1 Normal event dispatch
	dispatcher.Dispatch("user.registered", "User_888")
	assert.True(t, events.Tracker.Has("user.registered:User_888"))

	// 11.2 Until dispatch (Short-circuit intercept)
	// Low amount: passed (no return)
	res := dispatcher.Until("order.creating", float64(500))
	assert.Nil(t, res)

	// High amount: intercepted
	res = dispatcher.Until("order.creating", float64(50000))
	assert.Equal(t, "Amount exceeds risk threshold", res)
}

// 12. Test Console Command Handling
func TestConsoleKernel(t *testing.T) {
	app := bootstrap.BootApp()

	// route:list 输出路由表并成功返回。
	status := app.HandleCommand("route:list")
	assert.Equal(t, 0, status)

	// 未知命令返回非零退出码。
	status = app.HandleCommand("test:command")
	assert.NotEqual(t, 0, status)
}
