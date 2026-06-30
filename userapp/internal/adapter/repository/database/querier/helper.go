package querier

import "fmt"

// TenantTableSchema returns fully qualified table name with tenant schema prefix.
// The schema is double-quoted to safely handle UUID values containing hyphens.
// Example: TenantTableSchema("abc-123-uuid", "pocket_items") → `"abc-123-uuid".pocket_items`
func TenantTableSchema(schema, table string) string {
	return fmt.Sprintf(`"%s".%s`, schema, table)
}
