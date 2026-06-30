package seed

import "errors"

var (
	ErrSchemaNotFound                   = errors.New("schema not found")
	ErrSeedNotFound                     = errors.New("seed not found")
	ErrDeleteSeedData                   = errors.New("error delete seed data")
	ErrDeleteSeedDataBudgetPlanningItem = errors.New("error delete seed data - budget planning item")
)
