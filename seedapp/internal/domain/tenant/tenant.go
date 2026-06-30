package tenant

type Tenant struct {
	TenantCode string `json:"tenant_code" db:"tenant_code"`
	TenantName string `json:"tenant_name" db:"tenant_name"`
}

type ResultMigrate struct {
	Success int
	Fail    int
}
