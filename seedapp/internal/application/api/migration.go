package api

import (
	"seedapp/internal/application/model"
	"seedapp/internal/application/service"
	"seedapp/internal/pkg/formatter"
	"seedapp/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type MigrationAPI interface {
	RunMigration(*fiber.Ctx) error
}

type MigrationHandler struct {
	Service   service.MigrationService `inject:"migrationService"`
	Validator validator.Validator      `inject:"validator"`
}

func (h *MigrationHandler) RunMigration(c *fiber.Ctx) error {
	var payload model.MigrationRequest

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	if err := h.Validator.Validate(c.Context(), payload); err != nil {
		return err
	}

	result, err := h.Service.RunMigrations(c.Context(), &payload)
	if err != nil {
		return err
	}

	if len(result.Failed) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(formatter.NewErrorResponseList(formatter.DataConflict, "Failed to migrate schemas", c.Locals("traceId").(string), result.Failed))
	}

	return c.Status(fiber.StatusOK).JSON(formatter.NewSuccessResponse(formatter.Success, result))
}
