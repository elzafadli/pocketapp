package api

import (
	"seedapp/internal/application/model"
	"seedapp/internal/application/service"
	"seedapp/internal/pkg/formatter"
	"seedapp/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type SeedAPI interface {
	RunSeed(*fiber.Ctx) error
}

type SeedHandler struct {
	Service   service.SeedService `inject:"seedService"`
	Validator validator.Validator `inject:"validator"`
}

func (h *SeedHandler) RunSeed(c *fiber.Ctx) error {
	var payload model.SeedRequest

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	if err := h.Validator.Validate(c.Context(), payload); err != nil {
		return err
	}

	result, err := h.Service.RunSeeds(c.Context(), &payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(formatter.NewSuccessResponse(formatter.Success, result))
}
