package api

import (
	"userapp/internal/application/service"
	"userapp/internal/domain/shared/identity"
	"userapp/internal/domain/template"
	"userapp/internal/pkg/formatter"
	"userapp/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type TemplateAPI interface {
	Create(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Get(*fiber.Ctx) error
	Detail(*fiber.Ctx) error
}

type TemplateHandler struct {
	Service   service.TemplateService `inject:"templateService"`
	Validator validator.Validator     `inject:"validator"`
}

func (h *TemplateHandler) Create(c *fiber.Ctx) error {
	var payload *template.CreateTemplateRequest

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

func (h *TemplateHandler) Update(c *fiber.Ctx) error {
	var payload *template.UpdateTemplateRequest

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	payload.ID = identity.FromStringOrNil(c.Params("id"))

	res, err := h.Service.Update(c.Context(), payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(formatter.NewSuccessResponse(formatter.Success, res))
}

func (h *TemplateHandler) Get(c *fiber.Ctx) error {
	panic("not implemented") // TODO: Implement
}

func (h *TemplateHandler) Detail(c *fiber.Ctx) error {
	id := c.Params("id")
	data, err := h.Service.Find(c.Context(), identity.FromStringOrNil(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(formatter.NewSuccessResponse(formatter.Success, data))
}
