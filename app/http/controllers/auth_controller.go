package controllers

import (
	"github.com/go-think/flow"
	"github.com/go-think/think/view"
)

// AuthController handles authentication routes and views
type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// ShowLoginForm displays the login page
func (c *AuthController) ShowLoginForm(req *flow.Request) *flow.Response {
	if s := req.Session(); s != nil && s.Has("user_id") {
		return flow.Redirect("/dashboard")
	}
	return view.HTMLWithLayout("layouts/guest", "auth/login", map[string]any{
		"Title": "Sign In",
	}, req)
}

// Login processes login credentials
func (c *AuthController) Login(req *flow.Request, form *LoginRequest) *flow.Response {
	// Demo validation logic: accept password "password" or any >= 6 chars for non-empty user
	if form.Password != "password" && form.Password != "secret" && len(form.Password) < 6 {
		return flow.Redirect("/login").
			WithInput("email").
			WithErrors(map[string][]string{
				"general": {"Invalid credentials. Demo password is 'password'."},
			})
	}

	// Login user and establish session
	s := req.Session()
	if s != nil {
		s.Set("user_id", 101)
		s.Set("user_name", "Artisan Developer")
		s.Set("user_email", form.Email)
		s.Regenerate()
	}

	return flow.Redirect("/dashboard")
}

// ShowRegisterForm displays the registration page
func (c *AuthController) ShowRegisterForm(req *flow.Request) *flow.Response {
	if s := req.Session(); s != nil && s.Has("user_id") {
		return flow.Redirect("/dashboard")
	}
	return view.HTMLWithLayout("layouts/guest", "auth/register", map[string]any{
		"Title": "Register",
	}, req)
}

// Register processes user registration
func (c *AuthController) Register(req *flow.Request, form *RegisterRequest) *flow.Response {
	s := req.Session()
	if s != nil {
		s.Set("user_id", 102)
		s.Set("user_name", form.Name)
		s.Set("user_email", form.Email)
		s.Regenerate()
	}

	return flow.Redirect("/dashboard")
}

// Logout terminates authenticated session
func (c *AuthController) Logout(req *flow.Request) *flow.Response {
	if s := req.Session(); s != nil {
		s.Forget("user_id", "user_name", "user_email")
		s.RegenerateToken()
	}
	return flow.Redirect("/login")
}

// Dashboard displays user dashboard
func (c *AuthController) Dashboard(req *flow.Request) *flow.Response {
	return view.HTMLWithLayout("layouts/app", "dashboard", map[string]any{
		"Title": "Dashboard",
	}, req)
}
