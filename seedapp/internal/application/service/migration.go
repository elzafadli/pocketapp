package service

import (
	"context"
	"errors"
	"time"

	"seedapp/config"
	"seedapp/internal/application/model"
	"seedapp/internal/domain/migration"
	"seedapp/internal/domain/shared/entity"
	"seedapp/internal/domain/shared/identity"
	"seedapp/internal/domain/tenant"
)

//go:generate mockgen -destination=mocks/migration.go -package=mocks -source=migration.go
type MigrationService interface {
	RunDefaultMigration(ctx context.Context) error
	RunMigrations(ctx context.Context, request *model.MigrationRequest) (*model.MigrationResponse, error)
	RunTenantsMigration(ctx context.Context) error
}

type Migration struct {
	MigrationRepo migration.Repository    `inject:"migrationRepository"`
	TenantRepo    tenant.TenantRepository `inject:"tenantRepository"`
	Config        *config.Config          `inject:"config"`
}

func (s *Migration) RunDefaultMigration(ctx context.Context) error {
	// run main migrations
	_, err := s.MigrationRepo.RunMigrations(ctx, migration.MIGRATION_TYPE_MAIN, migration.SCHEMA_MAIN)
	if err != nil {
		return err
	}

	startedAt := time.Now()
	// run tenant migrations
	totalApplied, err := s.MigrationRepo.RunMigrations(ctx, migration.MIGRATION_TYPE_TENANT, migration.SCHEMA_DEFAULT)
	if err != nil {
		return err
	}

	// if there is no migration applied, and we already have a successful migration record, we don't need to create a new migration record
	if totalApplied == 0 {
		existingVersion, err := s.MigrationRepo.GetLatestMigrationVersion(ctx, migration.SCHEMA_DEFAULT)
		if err == nil && existingVersion != "" {
			return nil
		}
	}

	latestDefaultVersion := identity.NewID().String()[:8]

	migrationRecord := &migration.Migration{
		Entity:     entity.NewEntity(),
		Schema:     migration.SCHEMA_DEFAULT,
		Version:    latestDefaultVersion,
		Status:     migration.MIGRATION_STATUS_SUCCESS.String(),
		StartedAt:  startedAt,
		FinishedAt: time.Now(),
	}

	if err := s.MigrationRepo.CreateMigrationHistory(ctx, migrationRecord); err != nil {
		return err
	}

	return nil
}

func (s *Migration) RunMigrations(ctx context.Context, request *model.MigrationRequest) (*model.MigrationResponse, error) {
	latestDefaultVersion, err := s.MigrationRepo.GetLatestMigrationVersion(ctx, migration.SCHEMA_DEFAULT)
	if err != nil && !errors.Is(err, migration.ErrMigrationNotFound) {
		return nil, err
	}

	schemas, err := s.getTargetSchemas(ctx, request.Schemas)
	if err != nil {
		return nil, err
	}

	response := &model.MigrationResponse{
		Success: []string{},
		Failed:  []model.MigrationError{},
	}

	for _, sch := range schemas {
		if err := s.processSingleSchema(ctx, sch, latestDefaultVersion, response, request.TenantName); err != nil {
			return nil, err
		}
	}
	return response, nil
}

func (s *Migration) getTargetSchemas(ctx context.Context, requestSchemas []string) ([]string, error) {
	if len(requestSchemas) > 0 {
		return requestSchemas, nil
	}

	return s.MigrationRepo.GetSchemas(ctx)
}

func (s *Migration) processSingleSchema(ctx context.Context, schema string, latestDefaultVersion string, response *model.MigrationResponse, tenantName string) error {
	// Check if schema is not main or public, check if tenant exists, if not, create it
	if schema != migration.SCHEMA_MAIN && schema != migration.SCHEMA_DEFAULT {
		_, err := s.TenantRepo.Get(ctx, schema)
		if err != nil {
			tName := tenantName
			if tName == "" {
				tName = "Tenant " + schema
			}
			newTenant := &tenant.Tenant{
				TenantCode: schema,
				TenantName: tName,
			}
			_ = s.TenantRepo.Create(ctx, newTenant)
		}
	}

	latestVersion, err := s.MigrationRepo.GetLatestMigrationVersion(ctx, schema)
	if err != nil && !errors.Is(err, migration.ErrMigrationNotFound) {
		return err
	}

	if latestVersion != "" && latestVersion == latestDefaultVersion {
		return nil
	}

	latestStatus, err := s.MigrationRepo.GetLatestMigrationStatus(ctx, schema)
	if err != nil && !errors.Is(err, migration.ErrMigrationNotFound) {
		return err
	}

	if latestStatus == migration.MIGRATION_STATUS_RUNNING.String() {
		return nil
	}

	migrationRecord := &migration.Migration{
		Entity:    entity.NewEntity(),
		Schema:    schema,
		Version:   latestDefaultVersion,
		Status:    migration.MIGRATION_STATUS_RUNNING.String(),
		StartedAt: time.Now(),
	}

	if err := s.MigrationRepo.CreateMigrationHistory(ctx, migrationRecord); err != nil {
		return err
	}

	totalApplied, err := s.MigrationRepo.RunMigrations(ctx, migration.MIGRATION_TYPE_TENANT, schema)
	if err != nil {
		migrationRecord.Status = migration.MIGRATION_STATUS_FAILED.String()
		migrationRecord.Error = err.Error()
		response.Failed = append(response.Failed, model.MigrationError{
			Schema: schema,
			Error:  err.Error(),
		})
	} else {
		migrationRecord.Status = migration.MIGRATION_STATUS_SUCCESS.String()
		if totalApplied > 0 {
			response.Success = append(response.Success, schema)
		}
	}
	migrationRecord.FinishedAt = time.Now()

	if err := s.MigrationRepo.UpdateMigrationHistory(ctx, migrationRecord); err != nil {
		return err
	}

	return nil
}
func (s *Migration) RunTenantsMigration(ctx context.Context) error {

	if !s.Config.MigrateTenantStartup { // handle local development
		return nil
	}

	tenants, err := s.TenantRepo.GetAll(ctx, map[string]any{})
	if err != nil {
		return err
	}

	latestDefaultVersion, err := s.MigrationRepo.GetLatestMigrationVersion(ctx, migration.SCHEMA_DEFAULT)
	if err != nil && !errors.Is(err, migration.ErrMigrationNotFound) {
		return err
	}

	response := &model.MigrationResponse{
		Success: []string{},
		Failed:  []model.MigrationError{},
	}

	for _, tenant := range tenants {
		latestVersion, err := s.MigrationRepo.GetLatestMigrationVersion(ctx, tenant.TenantCode)
		if err != nil && !errors.Is(err, migration.ErrMigrationNotFound) {
			return err
		}

		if latestVersion != latestDefaultVersion {
			if err := s.processSingleSchema(ctx, tenant.TenantCode, latestDefaultVersion, response, tenant.TenantName); err != nil {
				return err
			}
		}
	}

	return nil
}
