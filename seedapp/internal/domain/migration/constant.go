package migration

type MigrationStatus string

const (
	MIGRATION_STATUS_RUNNING MigrationStatus = "running"
	MIGRATION_STATUS_SUCCESS MigrationStatus = "success"
	MIGRATION_STATUS_FAILED  MigrationStatus = "failed"
)

func (s MigrationStatus) String() string {
	return string(s)
}

type MigrationType string

const (
	MIGRATION_TYPE_MAIN   MigrationType = "main"
	MIGRATION_TYPE_TENANT MigrationType = "tenant"
)

func (s MigrationType) String() string {
	return string(s)
}

const (
	SCHEMA_MAIN    = "main"
	SCHEMA_DEFAULT = "public"
)
