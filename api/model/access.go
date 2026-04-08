package model

// CheckAccessResponse 权限查询响应
type CheckAccessResponse struct {
	HasAccess bool `json:"has_access"`
}

// ErrorResponse 通用错误响应
type ErrorResponse struct {
	Error string `json:"error"`
}
