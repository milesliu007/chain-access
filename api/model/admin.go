package model

// UserBalance represents a user's ERC20 token balance record
type UserBalance struct {
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Chain     string `json:"chain" gorm:"column:chain"`
	Address   string `json:"address" gorm:"column:address;index"`
	Token     string `json:"token" gorm:"column:token"`
	Balance   string `json:"balance" gorm:"column:balance"`
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (UserBalance) TableName() string { return "user_balances" }

// BalanceListRequest query params for balance list
type BalanceListRequest struct {
	Page    int    `form:"page"`
	Size    int    `form:"size"`
	Address string `form:"address"`
}

// BalanceListResponse paginated balance list response
type BalanceListResponse struct {
	Data  []UserBalance `json:"data"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

// AdminLoginRequest admin login via ERC721
type AdminLoginRequest struct {
	Address     string `json:"address" binding:"required"`
	Signature   string `json:"signature" binding:"required"`
	ChainID     string `json:"chain_id" binding:"required"`
	NFTContract string `json:"nft_contract" binding:"required"`
}

// AdminLoginResponse admin JWT response
type AdminLoginResponse struct {
	Token string `json:"token"`
}
