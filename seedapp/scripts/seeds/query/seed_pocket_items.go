package query

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/runsystemid/golog"
)

var seedPocketItemData = `
	INSERT INTO "%s".pocket_items (
		id,
		title,
		url,
		description,
		content_type,
		status,
		is_favorite,
		tags,
		created_at,
		updated_at,
		archived_at
	) VALUES
	(
		'10000000-0000-0000-0000-000000000001',
		'React Performance Guide',
		'https://example.com/react-performance',
		'A guide about React rendering optimization',
		'article',
		'unread',
		TRUE,
		'["frontend", "react"]'::jsonb,
		'2026-06-26T10:00:00Z',
		'2026-06-26T10:00:00Z',
		NULL
	),
	(
		'10000000-0000-0000-0000-000000000002',
		'Understanding TypeScript Generics',
		'https://example.com/typescript-generics',
		'Deep dive into TypeScript generics',
		'article',
		'reading',
		FALSE,
		'["typescript", "frontend"]'::jsonb,
		'2026-06-25T09:30:00Z',
		'2026-06-25T09:30:00Z',
		NULL
	),
	(
		'10000000-0000-0000-0000-000000000003',
		'Frontend System Design Notes',
		NULL,
		'Personal notes about scalable frontend architecture',
		'note',
		'read',
		TRUE,
		'["architecture", "frontend"]'::jsonb,
		'2026-06-24T08:00:00Z',
		'2026-06-24T08:00:00Z',
		NULL
	),
	(
		'10000000-0000-0000-0000-000000000004',
		'Archived Design Reference',
		'https://example.com/design-reference',
		'Old design reference kept for archive testing',
		'document',
		'archived',
		FALSE,
		'["design", "reference"]'::jsonb,
		'2026-06-20T08:00:00Z',
		'2026-06-28T08:00:00Z',
		'2026-06-28T08:00:00Z'
	)
	ON CONFLICT (id) DO NOTHING;
`

func SeedPocketItem(ctx context.Context, schema string, db *sqlx.DB) error {
	query := fmt.Sprintf(seedPocketItemData, schema)
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		golog.Error(ctx, "error seeding pocket items", err)
		return err
	}
	return nil
}
