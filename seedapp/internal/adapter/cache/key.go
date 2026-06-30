package cache

func KeyIgrisUserByCode(code string) string {
	return "igris:user:" + code
}

func KeyUserTenantByTenantCode(tenantCode string) string {
	return "user:tenant:" + tenantCode
}
