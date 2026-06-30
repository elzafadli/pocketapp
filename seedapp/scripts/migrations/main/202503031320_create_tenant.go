package migrationmain

import (
	migrate "github.com/rubenv/sql-migrate"
)

func CreateTenant() *migrate.Migration {
	mig := migrate.Migration{
		Id: "202503031320",
		Up: []string{
			`
			CREATE TABLE IF NOT EXISTS main.tenants (
				tenant_code uuid PRIMARY KEY DEFAULT public.uuid_generate_v4(),
				tenant_name varchar(250) NOT NULL,
				created_by uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
				created_at timestamp NOT NULL DEFAULT NOW(),
				updated_by uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
				updated_at timestamp NOT NULL DEFAULT NOW()
			);
			`,
		},
		Down: []string{
			`
			DROP TABLE IF EXISTS main.tenants;
			`,
		},
	}

	return &mig
}
