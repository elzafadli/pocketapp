package migrationmain

import (
	migrate "github.com/rubenv/sql-migrate"
)

func CreateTableUsers() *migrate.Migration {
	mig := migrate.Migration{
		Id: "202502271030",
		Up: []string{
			`
			SET TIME ZONE 'Asia/Jakarta';
			`,
			`
			CREATE TABLE IF NOT EXISTS users (
				user_code uuid PRIMARY KEY DEFAULT public.uuid_generate_v4(),
				user_name varchar(100) NOT NULL UNIQUE,
				tenant_default uuid NULL DEFAULT NULL,
				active_indicator varchar(1) NOT NULL DEFAULT 'Y',
				email varchar(100) NULL DEFAULT NULL,
				password varchar(255) NOT NULL,
				password_last_updated_at varchar(12) NULL DEFAULT NULL,
				expired_date varchar(12) NULL DEFAULT NULL,
				access_token varchar(255) DEFAULT NULL,
				created_by uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
				created_at timestamp NOT NULL DEFAULT NOW(),
				updated_by uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
				updated_at timestamp NOT NULL DEFAULT NOW()
			);
			INSERT INTO users (user_code, user_name, tenant_default, active_indicator, email, password, password_last_updated_at, expired_date, access_token, created_by, created_at, updated_by, updated_at) VALUES ('00000000-0000-0000-0000-000000000000', 'system', '00000000-0000-0000-0000-000000000000', 'Y', 'system@runsystem.id', 'system', NULL, NULL, NULL, '00000000-0000-0000-0000-000000000000', NOW(), '00000000-0000-0000-0000-000000000000', NOW());
			`,
		},
		Down: []string{
			`
			DROP TABLE IF EXISTS users;
			`,
		},
	}

	return &mig
}
