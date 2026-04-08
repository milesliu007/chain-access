package model

// ChallengeRequest 请求 challenge 的参数
type ChallengeRequest struct {
	Address string `json:"address" binding:"required"`
}

// ChallengeResponse challenge 响应
type ChallengeResponse struct {
	Challenge string `json:"challenge"`
}

// VerifyRequest 验证签名的参数
type VerifyRequest struct {
	Address   string `json:"address" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

// VerifyResponse 验证成功响应
type VerifyResponse struct {
	Token string `json:"token"`
}
