package auth

import (
	"userapp/internal/domain/user"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	TimeZone string `json:"time_zone"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl"`
}

type LoginResponse struct {
	Token      string       `json:"token"`
	TenantCode string       `json:"tenant_code"`
	User       UserResponse `json:"user"`
}

type LogoutRequest struct {
	TenantCode string `json:"tenant_code"`
	UserCode   string `json:"user_code"`
	SessionId  string `json:"session_id"`
}

type RegisterRequest struct {
	Name       string `json:"name" validate:"required,max=255"`
	Email      string `json:"email" validate:"required,email,max=255"`
	Password   string `json:"password" validate:"required,min=6"`
	TenantCode string `json:"tenant_code" validate:"omitempty,uuid"`
	TenantName string `json:"tenant_name" validate:"required,max=255"`
}

type RegisterResponse struct {
	User       *user.User `json:"user"`
	TenantCode string     `json:"tenant_code"`
	TenantName string     `json:"tenant_name"`
}
