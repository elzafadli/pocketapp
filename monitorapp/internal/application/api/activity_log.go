package api

import (
	"monitorapp/internal/application/service"
	"monitorapp/internal/domain/activity_log"
	"monitorapp/internal/pkg/formatter"
	"monitorapp/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type ActivityLogAPI interface {
	Create(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Delete(*fiber.Ctx) error
	List(*fiber.Ctx) error
	Detail(*fiber.Ctx) error
}

type ActivityLogHandler struct {
	Service   service.ActivityLogService `inject:"activityLogService"`
	Validator validator.Validator        `inject:"validator"`
}

func ActivityLogRoutes(router fiber.Router, authProtection fiber.Handler, h ActivityLogAPI) {
	router.Get("/activity-log", authProtection, h.List)
	router.Post("/activity-log", authProtection, h.Create)
	router.Put("/activity-log/:id", authProtection, h.Update)
	router.Delete("/activity-log/:id", authProtection, h.Delete)
	router.Get("/activity-log/:id", authProtection, h.Detail)
}

type ListActivityLogResponse struct {
	ActivityLogs []*activity_log.ActivityLog `json:"activity_logs"`
	Metadata     map[string]interface{}      `json:"metadata"`
}

func (h *ActivityLogHandler) Create(c *fiber.Ctx) error {
	var payload *activity_log.CreateActivityLogRequest

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	if err := h.Validator.Validate(c.Context(), payload); err != nil {
		return err
	}

	res, err := h.Service.Create(c.Context(), payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(formatter.NewSuccessResponse(formatter.Success, res))
}

func (h *ActivityLogHandler) Update(c *fiber.Ctx) error {
	var payload *activity_log.UpdateActivityLogRequest

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}
	payload.ID = id

	if err := h.Validator.Validate(c.Context(), payload); err != nil {
		return err
	}

	res, err := h.Service.Update(c.Context(), payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(formatter.NewSuccessResponse(formatter.Success, res))
}

func (h *ActivityLogHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	err := h.Service.Delete(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(formatter.NewSuccessResponse(formatter.Success, "Activity log successfully deleted"))
}

func (h *ActivityLogHandler) Detail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	data, err := h.Service.Find(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(formatter.NewSuccessResponse(formatter.Success, data))
}

func (h *ActivityLogHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	objectName := c.Query("object_name", "")
	recordID := c.Query("record_id", "")
	action := c.Query("action", "")
	changedBy := c.Query("changed_by", "")

	filter := map[string]interface{}{
		"limit":       limit,
		"offset":      offset,
		"object_name": objectName,
		"record_id":   recordID,
		"action":      action,
		"changed_by":  changedBy,
	}

	logs, total, err := h.Service.List(c.Context(), filter)
	if err != nil {
		return err
	}

	var currentPage int = page
	var totalPages int = 0
	if limit > 0 {
		totalPages = int((total + uint64(limit) - 1) / uint64(limit))
	}

	metadata := map[string]interface{}{
		"total":        total,
		"limit":        limit,
		"offset":       offset,
		"current_page": currentPage,
		"total_pages":  totalPages,
	}

	resp := ListActivityLogResponse{
		ActivityLogs: logs,
		Metadata:     metadata,
	}

	return c.Status(fiber.StatusOK).JSON(formatter.NewSuccessResponse(formatter.Success, resp))
}
