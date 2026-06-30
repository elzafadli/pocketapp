package seeds

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"seedapp/internal/domain/seed"
	queryseed "seedapp/scripts/seeds/query"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SeedFunc func(ctx context.Context, tenantType seed.SeedTenantType, schema string, db *sqlx.DB) error

type SeederService interface {
	AllSeeds() map[string]SeedFunc
	AllSeedsName(allSeeds map[string]SeedFunc) []string
	HashSeed(seeds []string) string
}

type Seeder struct {
}

func (s *Seeder) AllSeeds() map[string]SeedFunc {
	return map[string]SeedFunc{
		"pocket_item": SeedPocketItem,
	}
}

func (s *Seeder) AllSeedsName(allSeeds map[string]SeedFunc) []string {
	allSeedsKeys := make([]string, 0, len(allSeeds))
	for key := range allSeeds {
		allSeedsKeys = append(allSeedsKeys, key)
	}
	sort.Strings(allSeedsKeys)

	return allSeedsKeys
}

func (s *Seeder) HashSeed(seeds []string) string {
	sort.Strings(seeds)
	hash := sha256.New()
	hash.Write([]byte(strings.Join(seeds, "")))
	return hex.EncodeToString(hash.Sum(nil))
}


func SeedPocketItem(ctx context.Context, tenantType seed.SeedTenantType, schema string, db *sqlx.DB) error {
	err := queryseed.SeedPocketItem(ctx, schema, db)
	if err != nil {
		return err
	}

	return nil
}
