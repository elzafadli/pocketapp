package querier

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func AddFilterWhereInString(sql sq.SelectBuilder, values []string, fieldTable string) sq.SelectBuilder {
	if len(values) > 0 {
		if len(values) == 1 && values[0] != "" {
			sql = sql.Where(sq.Eq{fieldTable: values[0]})
		} else {
			sql = sql.Where(sq.Eq{fieldTable: values})
		}
	}
	return sql
}

func AddFilterWhereNotInString(sql sq.SelectBuilder, values []string, fieldTable string) sq.SelectBuilder {
	if len(values) > 0 {
		if len(values) == 1 && values[0] != "" {
			sql = sql.Where(sq.NotEq{fieldTable: values[0]})
		} else {
			sql = sql.Where(sq.NotEq{fieldTable: values})
		}
	}
	return sql
}

func AddFilterStartDate(sql sq.SelectBuilder, value string, fieldTable string) sq.SelectBuilder {
	if len(value) > 0 {
		sql = sql.Where(
			sq.GtOrEq{fieldTable: value},
		)
	}
	return sql
}
func AddFilterEndDate(sql sq.SelectBuilder, value string, fieldTable string) sq.SelectBuilder {
	if len(value) > 0 {
		sql = sql.Where(
			sq.LtOrEq{fieldTable: value},
		)
	}
	return sql
}

func TenantTableSchema(schema, tableName string) string {
	return fmt.Sprintf(`"%s".%s`, schema, tableName)
}
