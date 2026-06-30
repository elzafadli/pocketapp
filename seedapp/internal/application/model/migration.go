package model

type MigrationRequest struct {
	Schemas    []string `json:"schemas"`
	TenantName string   `json:"tenant_name"`
}

type MigrationResponse struct {
	Success []string         `json:"success"`
	Failed  []MigrationError `json:"failed"`
}

type MigrationError struct {
	Schema string `json:"schema"`
	Error  string `json:"error"`
}

type MWSMigrateModel struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

type MWSData struct {
	Body   any `json:"body"`
	Header any `json:"header"`
}

type KokaiRequest struct {
	Message any `json:"message"`
}
