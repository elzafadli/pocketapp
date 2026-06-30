package rest

import (
	"seedapp/internal/domain/auth"
	"seedapp/internal/domain/migration"
	"seedapp/internal/domain/seed"
	"seedapp/internal/domain/tenant"
	"seedapp/internal/pkg/formatter"

	"github.com/gofiber/fiber/v2"
)

var CodeMap = map[error]formatter.Status{
	// migration
	migration.ErrSchemaNotFound:    formatter.DataNotFound,
	migration.ErrMigrationNotFound: formatter.DataNotFound,
	migration.ErrFailedMigrate:     formatter.InvalidRequest,

	// seed
	seed.ErrSchemaNotFound:                   formatter.DataNotFound,
	seed.ErrSeedNotFound:                     formatter.DataNotFound,
	seed.ErrDeleteSeedData:                   formatter.InvalidRequest,
	seed.ErrDeleteSeedDataBudgetPlanningItem: formatter.DataConflict,

	// auth
	auth.ErrInvalidBasicAuth: formatter.Unauthorized,

	// tenant
	tenant.ErrTenantNotFound: formatter.DataNotFound,
}

var StatusMap = map[error]int{
	// migration
	migration.ErrSchemaNotFound:    fiber.StatusNotFound,
	migration.ErrMigrationNotFound: fiber.StatusNotFound,
	migration.ErrFailedMigrate:     fiber.StatusBadRequest,

	// seed
	seed.ErrSchemaNotFound:                   fiber.StatusNotFound,
	seed.ErrSeedNotFound:                     fiber.StatusNotFound,
	seed.ErrDeleteSeedData:                   fiber.StatusBadRequest,
	seed.ErrDeleteSeedDataBudgetPlanningItem: fiber.ErrBadRequest.Code,

	// auth
	auth.ErrInvalidBasicAuth: fiber.StatusUnauthorized,

	// tenant
	tenant.ErrTenantNotFound: fiber.StatusNotFound,
}
