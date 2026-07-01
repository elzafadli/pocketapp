package api

import (
	"net/url"
	"strings"

	"userapp/internal/application/service"
	"userapp/internal/domain/pocket"
	"userapp/internal/domain/shared"
	"userapp/internal/domain/shared/identity"
	"userapp/internal/pkg/custommiddleware"
	"userapp/internal/pkg/formatter"
	"userapp/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type PocketAPI interface {
	Create(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Get(*fiber.Ctx) error
	Detail(*fiber.Ctx) error
	Delete(*fiber.Ctx) error
	UpdateStatus(*fiber.Ctx) error
	ToggleFavorite(*fiber.Ctx) error
	Summarize(*fiber.Ctx) error
	GetSummary(*fiber.Ctx) error
}

type PocketHandler struct {
	Service   service.PocketService `inject:"pocketService"`
	Validator validator.Validator   `inject:"validator"`
}

func PocketRoutes(router fiber.Router, authMiddleware custommiddleware.AuthMiddlewareService, h PocketAPI) {
	router.Get("/pockets", authMiddleware.JwtProtection(), h.Get)
	router.Get("/pockets/summary", authMiddleware.JwtProtection(), h.GetSummary)
	router.Post("/pockets", authMiddleware.JwtProtection(), h.Create)
	router.Get("/pockets/:id", authMiddleware.JwtProtection(), h.Detail)
	router.Put("/pockets/:id", authMiddleware.JwtProtection(), h.Update)
	router.Delete("/pockets/:id", authMiddleware.JwtProtection(), h.Delete)
	router.Patch("/pockets/:id/status", authMiddleware.JwtProtection(), h.UpdateStatus)
	router.Patch("/pockets/:id/favorite", authMiddleware.JwtProtection(), h.ToggleFavorite)
	router.Post("/pockets/:id/summarize", authMiddleware.JwtProtection(), h.Summarize)
}

// validateURL handles URL validation separately because the requirement is conditional:
// URL is required for article/video/document, but optional for note.
func validateURL(contentType, rawUrl string) *shared.ValidationErrorDetail {
	urlTrimmed := strings.TrimSpace(rawUrl)
	contentTypeTrimmed := strings.TrimSpace(contentType)
	isUrlRequired := contentTypeTrimmed == "article" || contentTypeTrimmed == "video" || contentTypeTrimmed == "document"

	if isUrlRequired && urlTrimmed == "" {
		return &shared.ValidationErrorDetail{Field: "url", Message: "URL is required"}
	}
	if urlTrimmed != "" {
		_, err := url.ParseRequestURI(urlTrimmed)
		if err != nil || !strings.Contains(urlTrimmed, ".") {
			return &shared.ValidationErrorDetail{Field: "url", Message: "URL is invalid"}
		}
	}
	return nil
}

// validatorErrorsToDetails converts ErrorMap from the validator package into shared.ValidationErrorDetail slice.
func validatorErrorsToDetails(err error) []shared.ValidationErrorDetail {
	errMap, ok := err.(*validator.ErrorMap)
	if !ok {
		return []shared.ValidationErrorDetail{{Field: "_", Message: err.Error()}}
	}
	var details []shared.ValidationErrorDetail
	for field, e := range errMap.Errors {
		details = append(details, shared.ValidationErrorDetail{Field: field, Message: e.Error()})
	}
	return details
}

func (h *PocketHandler) Create(c *fiber.Ctx) error {
	schema := tenantSchema(c)

	var payload pocket.CreatePocketRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(formatter.NewBadRequestResponse(err.Error()))
	}

	// Validate via struct tags (title, contentType, description, tags)
	var validationErrors []shared.ValidationErrorDetail
	if err := h.Validator.Validate(c.Context(), payload); err != nil {
		validationErrors = append(validationErrors, validatorErrorsToDetails(err)...)
	}

	// Custom URL validation (conditional based on contentType)
	if urlErr := validateURL(payload.ContentType, payload.URL); urlErr != nil {
		validationErrors = append(validationErrors, *urlErr)
	}

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(formatter.NewValidationErrorResponse(validationErrors))
	}

	res, err := h.Service.Create(c.Context(), schema, &payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(formatter.NewInternalErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    res,
		"message": "Pocket item created successfully",
	})
}

func (h *PocketHandler) Update(c *fiber.Ctx) error {
	schema := tenantSchema(c)

	var payload pocket.UpdatePocketRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(formatter.NewBadRequestResponse(err.Error()))
	}

	idStr := c.Params("id")
	payload.ID = identity.FromStringOrNil(idStr)
	if payload.ID.UUID == identity.NewZeroID().UUID {
		return c.Status(fiber.StatusNotFound).JSON(formatter.NewNotFoundResponse(pocket.ErrCodePocketNotFound, pocket.ErrMsgPocketNotFound))
	}

	// Validate via struct tags (title, contentType, description, tags)
	var validationErrors []shared.ValidationErrorDetail
	if err := h.Validator.Validate(c.Context(), payload); err != nil {
		validationErrors = append(validationErrors, validatorErrorsToDetails(err)...)
	}

	// Custom URL validation (conditional based on contentType)
	if urlErr := validateURL(payload.ContentType, payload.URL); urlErr != nil {
		validationErrors = append(validationErrors, *urlErr)
	}

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(formatter.NewValidationErrorResponse(validationErrors))
	}

	res, err := h.Service.Update(c.Context(), schema, &payload)
	if err != nil {
		if err == pocket.ErrPocketNotFound {
			return c.Status(fiber.StatusNotFound).JSON(formatter.NewNotFoundResponse(pocket.ErrCodePocketNotFound, pocket.ErrMsgPocketNotFound))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(formatter.NewInternalErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    res,
		"message": "Pocket item updated successfully",
	})
}

func (h *PocketHandler) Get(c *fiber.Ctx) error {
	schema := tenantSchema(c)

	query := &pocket.PocketListQuery{
		Search:      c.Query("search", ""),
		Status:      c.Query("status", ""),
		ContentType: c.Query("type", ""),
		Page:        c.QueryInt("page", 1),
		Limit:       c.QueryInt("limit", 10),
		Sort:        c.Query("sort", "createdAt:desc"),
	}

	if favoriteStr := c.Query("favorite", ""); favoriteStr == "true" {
		v := true
		query.Favorite = &v
	} else if favoriteStr == "false" {
		v := false
		query.Favorite = &v
	}

	list, total, err := h.Service.List(c.Context(), schema, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(formatter.NewInternalErrorResponse(err.Error()))
	}

	totalPage := int((total + uint64(query.Limit) - 1) / uint64(query.Limit))
	if total == 0 {
		totalPage = 0
	}

	return c.Status(fiber.StatusOK).JSON(pocket.PocketListResponse{
		Data: list,
		Meta: shared.MetaResponse{
			Page:      query.Page,
			Limit:     query.Limit,
			Total:     int64(total),
			TotalPage: totalPage,
		},
	})
}

func (h *PocketHandler) Detail(c *fiber.Ctx) error {
	schema := tenantSchema(c)

	idStr := c.Params("id")
	id := identity.FromStringOrNil(idStr)
	if id.UUID == identity.NewZeroID().UUID {
		return c.Status(fiber.StatusNotFound).JSON(formatter.NewNotFoundResponse(pocket.ErrCodePocketNotFound, pocket.ErrMsgPocketNotFound))
	}

	res, err := h.Service.Find(c.Context(), schema, id)
	if err != nil {
		if err == pocket.ErrPocketNotFound {
			return c.Status(fiber.StatusNotFound).JSON(formatter.NewNotFoundResponse(pocket.ErrCodePocketNotFound, pocket.ErrMsgPocketNotFound))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(formatter.NewInternalErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}

func (h *PocketHandler) Delete(c *fiber.Ctx) error {
	schema := tenantSchema(c)

	idStr := c.Params("id")
	id := identity.FromStringOrNil(idStr)
	if id.UUID == identity.NewZeroID().UUID {
		return c.Status(fiber.StatusNotFound).JSON(formatter.NewNotFoundResponse(pocket.ErrCodePocketNotFound, pocket.ErrMsgPocketNotFound))
	}

	err := h.Service.Delete(c.Context(), schema, id)
	if err != nil {
		if err == pocket.ErrPocketNotFound {
			return c.Status(fiber.StatusNotFound).JSON(formatter.NewNotFoundResponse(pocket.ErrCodePocketNotFound, pocket.ErrMsgPocketNotFound))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(formatter.NewInternalErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"id": id.String(),
		},
		"message": "Pocket item archived successfully",
	})
}

func (h *PocketHandler) UpdateStatus(c *fiber.Ctx) error {
	schema := tenantSchema(c)

	var payload pocket.UpdateStatusRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(formatter.NewBadRequestResponse(err.Error()))
	}

	idStr := c.Params("id")
	id := identity.FromStringOrNil(idStr)
	if id.UUID == identity.NewZeroID().UUID {
		return c.Status(fiber.StatusNotFound).JSON(formatter.NewNotFoundResponse(pocket.ErrCodePocketNotFound, pocket.ErrMsgPocketNotFound))
	}

	statusTrimmed := strings.TrimSpace(payload.Status)
	if statusTrimmed != "unread" && statusTrimmed != "reading" && statusTrimmed != "read" && statusTrimmed != "archived" {
		return c.Status(fiber.StatusBadRequest).JSON(formatter.NewValidationErrorResponse(
			[]shared.ValidationErrorDetail{
				{Field: "status", Message: "Status is invalid"},
			},
		))
	}

	res, err := h.Service.UpdateStatus(c.Context(), schema, id, statusTrimmed)
	if err != nil {
		if err == pocket.ErrPocketNotFound {
			return c.Status(fiber.StatusNotFound).JSON(formatter.NewNotFoundResponse(pocket.ErrCodePocketNotFound, pocket.ErrMsgPocketNotFound))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(formatter.NewInternalErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"id":        res.ID.String(),
			"status":    res.Status,
			"updatedAt": res.UpdatedAt,
		},
		"message": "Status updated successfully",
	})
}

func (h *PocketHandler) ToggleFavorite(c *fiber.Ctx) error {
	schema := tenantSchema(c)

	var payload pocket.ToggleFavoriteRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(formatter.NewBadRequestResponse(err.Error()))
	}

	idStr := c.Params("id")
	id := identity.FromStringOrNil(idStr)
	if id.UUID == identity.NewZeroID().UUID {
		return c.Status(fiber.StatusNotFound).JSON(formatter.NewNotFoundResponse(pocket.ErrCodePocketNotFound, pocket.ErrMsgPocketNotFound))
	}

	res, err := h.Service.ToggleFavorite(c.Context(), schema, id, payload.IsFavorite)
	if err != nil {
		if err == pocket.ErrPocketNotFound {
			return c.Status(fiber.StatusNotFound).JSON(formatter.NewNotFoundResponse(pocket.ErrCodePocketNotFound, pocket.ErrMsgPocketNotFound))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(formatter.NewInternalErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"id":         res.ID.String(),
			"isFavorite": res.IsFavorite,
			"updatedAt":  res.UpdatedAt,
		},
		"message": "Favorite updated successfully",
	})
}

func (h *PocketHandler) Summarize(c *fiber.Ctx) error {
	schema := tenantSchema(c)

	idStr := c.Params("id")
	id := identity.FromStringOrNil(idStr)
	if id.UUID == identity.NewZeroID().UUID {
		return c.Status(fiber.StatusNotFound).JSON(formatter.NewNotFoundResponse(pocket.ErrCodePocketNotFound, pocket.ErrMsgPocketNotFound))
	}

	summary, err := h.Service.Summarize(c.Context(), schema, id)
	if err != nil {
		if err == pocket.ErrPocketNotFound {
			return c.Status(fiber.StatusNotFound).JSON(formatter.NewNotFoundResponse(pocket.ErrCodePocketNotFound, pocket.ErrMsgPocketNotFound))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(formatter.NewInternalErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    summary,
		"message": "Summary generated successfully",
	})
}

func (h *PocketHandler) GetSummary(c *fiber.Ctx) error {
	schema := tenantSchema(c)

	res, err := h.Service.GetSummary(c.Context(), schema)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(formatter.NewInternalErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    res,
		"message": "Dashboard summary retrieved successfully",
	})
}
