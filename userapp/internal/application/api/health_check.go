package api

import (
	"os"
	"strings"

	"userapp/internal/adapter/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/runsystemid/gocache"
)

type HealthCheckAPI interface {
	Ping(*fiber.Ctx) error
	Ready(*fiber.Ctx) error
	Version(*fiber.Ctx) error
}

type HealthCheckHandler struct {
	Database repository.Sqlx `inject:"database"`
	Cache    gocache.Service `inject:"cache"`
}

func (h *HealthCheckHandler) Ping(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "pong"})
}

func (h *HealthCheckHandler) Ready(c *fiber.Ctx) error {
	message := make(map[string]string)

	err := h.Database.Ping()
	if err != nil {
		message["database"] = "not ready"
	} else {
		message["database"] = "ready"
	}

	err = h.Cache.Ping(c.Context())
	if err != nil {
		message["cache"] = "not ready"
	} else {
		message["cache"] = "ready"
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(message)
	}

	return c.Status(fiber.StatusOK).JSON(message)
}

func (h *HealthCheckHandler) Version(c *fiber.Ctx) error {
	data, err := os.ReadFile("version.txt")
	version := "unknown"

	if err == nil {
		version = strings.TrimSpace(string(data))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"version": version})
}
