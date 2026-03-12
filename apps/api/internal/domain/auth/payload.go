package auth

// UpdateProfilePayload is the request body for PUT /api/v1/auth/profile.
// It allows the authenticated user to update their own name and email.
type UpdateProfilePayload struct {
	Name  string `json:"name"  binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email"`
}

// ChangePasswordPayload is the request body for PUT /api/v1/auth/password.
// Only valid for users with auth_provider = "local".
// The current password is verified before the new one is stored.
type ChangePasswordPayload struct {
	CurrentPassword string `json:"current_password" binding:"required,min=1"`
	NewPassword     string `json:"new_password"     binding:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=72"`
}
