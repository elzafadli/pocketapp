package api

import (
	"userapp/internal/application/service"
	"userapp/internal/domain/auth"
	"userapp/internal/domain/shared"
	"userapp/internal/pkg/custommiddleware"
	"userapp/internal/pkg/formatter"
	"userapp/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthAPI interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
}

type AuthHandler struct {
	AuthService service.AuthService `inject:"authServiceLoggable"`
	Validator   validator.Validator `inject:"validator"`
}

func AuthRoutes(router fiber.Router, authMiddleware custommiddleware.AuthMiddlewareService, h AuthAPI) {
	router.Post("/login", authMiddleware.ApiKeyProtection(), h.Login)
	router.Post("/logout", authMiddleware.BasicAuthProtection(), h.Logout)
	router.Post("/register", authMiddleware.ApiKeyProtection(), h.Register)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req auth.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if len(req.TimeZone) == 0 {
		req.TimeZone = "Asia/Jakarta"
	}

	if err := h.Validator.Validate(c.Context(), req); err != nil {
		return err
	}

	res, err := h.AuthService.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(shared.ErrorResponse{
			Code:    "INVALID_CREDENTIAL",
			Message: "Email or password is incorrect",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    res,
		"message": "Login success",
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req auth.LogoutRequest

	// Extract context from Locals (Assuming token parses them to Locals)
	tenantCode, _ := c.Locals("tenant_code").(string)
	userCode, _ := c.Locals("user_code").(string)
	sessionId, _ := c.Locals("session_id").(string)

	req.TenantCode = tenantCode
	req.UserCode = userCode
	req.SessionId = sessionId

	if err := h.AuthService.Logout(c.Context(), &req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(formatter.NewSuccessResponse(formatter.Success, nil))
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req auth.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := h.Validator.Validate(c.Context(), req); err != nil {
		return err
	}

	res, err := h.AuthService.Register(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(formatter.NewSuccessResponse(formatter.Success, res))
}
