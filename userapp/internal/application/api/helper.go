package api

import "github.com/gofiber/fiber/v2"

// tenantSchema extracts the tenant schema from the JWT Locals set by the auth middleware.
func tenantSchema(c *fiber.Ctx) string {
	schema, _ := c.Locals("tenant_code").(string)
	return schema
}
