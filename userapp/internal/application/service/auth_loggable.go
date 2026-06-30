package service

import (
	"context"

	"userapp/internal/domain/auth"
)

type AuthLoggable struct {
	Next               AuthService        `inject:"authService"`
	ActivityLogService ActivityLogService `inject:"activityLogService"`
}

func (s *AuthLoggable) Login(ctx context.Context, req auth.LoginRequest) (*auth.LoginResponse, error) {
	resp, err := s.Next.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	go func() {
		// Mask the password for security
		logReq := req
		logReq.Password = "[MASKED]"

		_ = s.ActivityLogService.Send(context.Background(), &ActivityLogRequest{
			ObjectName: "User",
			RecordID:   req.Email,
			Action:     "LOGIN",
			ChangedBy:  req.Email,
			Request:    logReq,
			Response:   nil,
		})
	}()

	return resp, nil
}

func (s *AuthLoggable) Logout(ctx context.Context, req *auth.LogoutRequest) error {
	err := s.Next.Logout(ctx, req)
	if err != nil {
		return err
	}

	go func() {
		// Copy the request to avoid potential concurrency issues
		var logReq auth.LogoutRequest
		if req != nil {
			logReq = *req
		}

		_ = s.ActivityLogService.Send(context.Background(), &ActivityLogRequest{
			ObjectName: "User",
			RecordID:   logReq.UserCode,
			Action:     "LOGOUT",
			ChangedBy:  logReq.UserCode,
			Request:    logReq,
			Response:   map[string]interface{}{"status": "success"},
		})
	}()

	return nil
}

func (s *AuthLoggable) Register(ctx context.Context, req auth.RegisterRequest) (*auth.RegisterResponse, error) {
	resp, err := s.Next.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	go func() {
		// Mask the password for security
		logReq := req
		logReq.Password = "[MASKED]"

		_ = s.ActivityLogService.Send(context.Background(), &ActivityLogRequest{
			ObjectName: "User",
			RecordID:   req.Email,
			Action:     "REGISTER",
			ChangedBy:  req.Email,
			Request:    logReq,
			Response:   resp,
		})
	}()

	return resp, nil
}
