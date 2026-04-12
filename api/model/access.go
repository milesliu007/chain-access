package model

// CheckAccessRequest ERC-20 权限查询请求
type CheckAccessRequest struct {
	ChainID         string `json:"chain_id" binding:"required"`
	Address         string `json:"address"`
	ContractAddress string `json:"contract_address" binding:"required"`
}

// CheckAccessResponse 权限查询响应
type CheckAccessResponse struct {
	HasAccess bool `json:"has_access"`
}

// CheckNFTRequest NFT 查询请求
type CheckNFTRequest struct {
	ChainID         string `json:"chain_id" binding:"required"`
	Address         string `json:"address"`
	ContractAddress string `json:"contract_address" binding:"required"`
}

// CheckNFTResponse NFT 查询响应
type CheckNFTResponse struct {
	HasNFT   bool     `json:"has_nft"`
	TokenIDs []string `json:"token_ids"`
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
