package model

// CheckAccessResponse 权限查询响应
type CheckAccessResponse struct {
	HasAccess bool `json:"has_access"`
}

// ChainInfo 链信息（返回给前端）
type ChainInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	ChainID int64  `json:"chain_id"`
}

// ChainsResponse 链列表响应
type ChainsResponse struct {
	Chains []ChainInfo `json:"chains"`
}

// ErrorResponse 通用错误响应
type ErrorResponse struct {
	Error string `json:"error"`
}
