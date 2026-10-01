package controllers

// StoreUserRequest represents an incoming FormRequest for creating a user.
type StoreUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Rules defines the validation rules matching Laravel FormRequest.
func (r *StoreUserRequest) Rules() map[string]string {
	return map[string]string{
		"name":  "required|min:2",
		"email": "required|email",
	}
}

// LoginRequest represents an incoming request to authenticate a user.
type LoginRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	Remember bool   `form:"remember" json:"remember"`
}

// Rules defines the validation rules for LoginRequest.
func (r *LoginRequest) Rules() map[string]string {
	return map[string]string{
		"email":    "required|email",
		"password": "required|min:6",
	}
}

// RegisterRequest represents an incoming request to register a new user.
type RegisterRequest struct {
	Name                 string `form:"name" json:"name"`
	Email                string `form:"email" json:"email"`
	Password             string `form:"password" json:"password"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation"`
}

// Rules defines the validation rules for RegisterRequest.
func (r *RegisterRequest) Rules() map[string]string {
	return map[string]string{
		"name":     "required|min:2",
		"email":    "required|email",
		"password": "required|min:6|confirmed",
	}
}
