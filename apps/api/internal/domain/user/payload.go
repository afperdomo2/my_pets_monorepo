package user

// CreateUserPayload is the shape accepted from HTTP clients on user creation.
// Password is required because local accounts always need one.
type CreateUserPayload struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// UpdateUserPayload is the shape accepted from HTTP clients on user update.
// Password is not updatable via this endpoint.
type UpdateUserPayload struct {
	Name  string `json:"name"  binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email"`
}

// GoogleUserInfo contains the profile data returned by Google's userinfo endpoint.
type GoogleUserInfo struct {
	GoogleID string
	Email    string
	Name     string
}
