package config

// ChainConfig 链配置
type ChainConfig struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	RPCURL  string `json:"-"`
	ChainID int64  `json:"chain_id"`
}

// DefaultChains 返回内置链注册表，新增链只需在此追加
func DefaultChains(infuraKey string) []ChainConfig {
	return []ChainConfig{
		{
			ID:      "ethereum",
			Name:    "Ethereum Mainnet",
			RPCURL:  "https://mainnet.infura.io/v3/" + infuraKey,
			ChainID: 1,
		},
		{
			ID:      "avax-fuji",
			Name:    "Avalanche Fuji Testnet",
			RPCURL:  "https://avalanche-fuji.infura.io/v3/" + infuraKey,
			ChainID: 43113,
		},
	}
}
