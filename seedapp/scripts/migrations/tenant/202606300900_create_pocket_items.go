package migrations

import migrate "github.com/rubenv/sql-migrate"

func CreatePocketItems(schema string) *migrate.Migration {
	mig := migrate.Migration{
		Id: "202606300900",
		Up: []string{
			`
			CREATE TYPE "` + schema + `".pocket_content_type AS ENUM ('article', 'video', 'document', 'note');
			CREATE TYPE "` + schema + `".pocket_status AS ENUM ('unread', 'reading', 'read', 'archived');

			CREATE TABLE IF NOT EXISTS "` + schema + `".pocket_items (
				id UUID PRIMARY KEY,
				title VARCHAR(255) NOT NULL,
				url TEXT NULL,
				description TEXT NULL,
				content_type "` + schema + `".pocket_content_type NOT NULL,
				status "` + schema + `".pocket_status NOT NULL DEFAULT 'unread',
				is_favorite BOOLEAN NOT NULL DEFAULT FALSE,
				tags JSONB NOT NULL DEFAULT '[]'::jsonb,
				created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				archived_at TIMESTAMPTZ NULL,

				CONSTRAINT pocket_items_url_required_check
					CHECK (
						content_type::TEXT = 'note'
						OR (url IS NOT NULL AND length(trim(url)) > 0)
					)
			);

			CREATE INDEX IF NOT EXISTS idx_pocket_items_active_created_at
				ON "` + schema + `".pocket_items(created_at DESC)
				WHERE archived_at IS NULL;

			CREATE INDEX IF NOT EXISTS idx_pocket_items_archived_at
				ON "` + schema + `".pocket_items(archived_at DESC)
				WHERE archived_at IS NOT NULL;

			CREATE INDEX IF NOT EXISTS idx_pocket_items_status_active
				ON "` + schema + `".pocket_items(status, created_at DESC)
				WHERE archived_at IS NULL;

			CREATE INDEX IF NOT EXISTS idx_pocket_items_content_type_active
				ON "` + schema + `".pocket_items(content_type, created_at DESC)
				WHERE archived_at IS NULL;

			CREATE INDEX IF NOT EXISTS idx_pocket_items_favorite_active
				ON "` + schema + `".pocket_items(is_favorite, created_at DESC)
				WHERE archived_at IS NULL;

			CREATE INDEX IF NOT EXISTS idx_pocket_items_title_active
				ON "` + schema + `".pocket_items(lower(title))
				WHERE archived_at IS NULL;

			CREATE INDEX IF NOT EXISTS idx_pocket_items_tags
				ON "` + schema + `".pocket_items USING gin(tags);
			`,
		},
		Down: []string{
			`
			DROP TABLE IF EXISTS "` + schema + `".pocket_items CASCADE;
			DROP TYPE IF EXISTS "` + schema + `".pocket_content_type CASCADE;
			DROP TYPE IF EXISTS "` + schema + `".pocket_status CASCADE;
			`,
		},
	}

	return &mig
}
