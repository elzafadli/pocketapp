package migrationmain

import (
	migrate "github.com/rubenv/sql-migrate"
)

func CreateUserTenant() *migrate.Migration {
	mig := migrate.Migration{
		Id: "202503031315",
		Up: []string{
			`
			SET TIME ZONE 'Asia/Jakarta';
			`,
			`
			CREATE TABLE IF NOT EXISTS main.user_tenants (
				user_code uuid NOT NULL,
				tenant_code uuid NOT NULL,
				active_indicator varchar(1) NOT NULL DEFAULT 'Y',
				created_by uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
				created_at timestamp NOT NULL DEFAULT NOW(),
				updated_by uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
				updated_at timestamp NOT NULL DEFAULT NOW(),
				CONSTRAINT user_tenant_pk PRIMARY KEY (user_code,tenant_code)
			);
			`,
			`
			CREATE INDEX IF NOT EXISTS user_tenant_user_code_idx ON main.user_tenants (user_code,tenant_code);
			`,
		},
		Down: []string{
			`
			DROP TABLE IF EXISTS main.user_tenants;
			`,
		},
	}

	return &mig
}
