package service

import (
	"context"
	"errors"
	"time"

	"userapp/config"
	"userapp/internal/adapter/kaisel"
	"userapp/internal/domain/auth"
	"userapp/internal/domain/user"
	"userapp/internal/domain/user_tenant"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/runsystemid/gocache"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context, req *auth.LogoutRequest) error
	Register(ctx context.Context, req auth.RegisterRequest) (*auth.RegisterResponse, error)
}

type Auth struct {
	Conf           *config.Config          `inject:"config"`
	UserRepo       user.Repository         `inject:"userRepository"`
	UserTenantRepo user_tenant.Repository  `inject:"userTenantRepository"`
	Cache          gocache.Service         `inject:"cache"`
	Kaisel         kaisel.KaiselService    `inject:"kaisel"`
}

func (s *Auth) Login(ctx context.Context, req auth.LoginRequest) (*auth.LoginResponse, error) {
	u, err := s.UserRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	roles := []string{}

	secret := s.Conf.JwtSecret
	if secret == "" {
		secret = "default_secret_key"
	}

	now := time.Now()
	accessTokenExpiredAt := now.Add(24 * time.Hour)

	// Generate Access Token
	accessClaims := jwt.MapClaims{
		"email":       req.Email,
		"tenant_code": u.TenantDefault,
		"roles":       roles,
		"exp":         accessTokenExpiredAt.Unix(),
		"iat":         now.Unix(),
		"type":        "access",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	res := &auth.LoginResponse{
		Token:      accessTokenString,
		TenantCode: u.TenantDefault,
		User: auth.UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			AvatarURL: "https://example.com/avatar.png",
		},
	}

	return res, nil
}

func (s *Auth) Logout(ctx context.Context, req *auth.LogoutRequest) error {
	// TODO: Implement actual logout logic (e.g., invalidate token/session)
	return nil
}

func (s *Auth) Register(ctx context.Context, req auth.RegisterRequest) (*auth.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	tenantCode := req.TenantCode
	if tenantCode == "" {
		tenantCode = uuid.New().String()
	}

	u := &user.User{
		Name:            req.Name,
		TenantDefault:   tenantCode,
		ActiveIndicator: "Y",
		Email:           req.Email,
		Password:        string(hashedPassword),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	newUser, err := s.UserRepo.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	// Insert ke user_tenants agar user hanya bisa akses data tenant-nya sendiri
	ut := &user_tenant.UserTenant{
		UserCode:        newUser.ID,
		TenantCode:      tenantCode,
		ActiveIndicator: "Y",
		CreatedBy:       newUser.ID,
		CreatedAt:       now,
		UpdatedBy:       newUser.ID,
		UpdatedAt:       now,
	}
	if err := s.UserTenantRepo.Create(ctx, ut); err != nil {
		return nil, err
	}

	traceId, _ := ctx.Value("traceId").(string)
	migrateReq := &kaisel.MigrateRequest{
		Schemas:    []string{tenantCode},
		TenantName: req.TenantName,
		TraceId:    traceId,
	}

	err = s.Kaisel.Migrate(migrateReq)
	if err != nil {
		return nil, err
	}

	res := &auth.RegisterResponse{
		User:       newUser,
		TenantCode: tenantCode,
		TenantName: req.TenantName,
	}

	return res, nil
}
