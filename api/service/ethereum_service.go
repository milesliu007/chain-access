package service

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"sync"
	"time"

	"chain-access/api/config"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const rpcTimeout = 10 * time.Second

// EthereumService 链上查询服务接口
type EthereumService interface {
	// CheckERC20Balance 查询指定链上 ERC-20 余额是否 > 0
	CheckERC20Balance(chainID, walletAddress, contractAddress string) (bool, error)
	// CheckERC721Ownership 查询指定链上 ERC-721 NFT 持有情况，返回是否持有及 tokenId 列表
	CheckERC721Ownership(chainID, walletAddress, contractAddress string) (bool, []*big.Int, error)
	// CheckERC1155Balance 查询指定链上 ERC-1155 代币持有数量
	CheckERC1155Balance(chainID, walletAddress, contractAddress, tokenID string) (bool, *big.Int, error)
	// GetChains 返回支持的链列表
	GetChains() []config.ChainConfig
	// Close 关闭所有连接
	Close()
}

// chainClient 单条链的 ethclient
type chainClient struct {
	client *ethclient.Client
	config config.ChainConfig
}

// EthereumServiceImpl 多链查询服务实现
type EthereumServiceImpl struct {
	mu      sync.RWMutex
	clients map[string]*chainClient // key: chain ID
	chains  []config.ChainConfig
}

// NewEthereumService 创建多链查询服务
func NewEthereumService(chains []config.ChainConfig, proxyAddr string) (EthereumService, error) {
	svc := &EthereumServiceImpl{
		clients: make(map[string]*chainClient, len(chains)),
		chains:  chains,
	}

	for _, chain := range chains {
		client, err := dialChain(chain.RPCURL, proxyAddr)
		if err != nil {
			svc.Close()
			return nil, fmt.Errorf("failed to connect to chain %s: %w", chain.ID, err)
		}
		svc.clients[chain.ID] = &chainClient{client: client, config: chain}
	}

	return svc, nil
}

func dialChain(rpcURL, proxyAddr string) (*ethclient.Client, error) {
	var opts []rpc.ClientOption

	if proxyAddr != "" {
		proxyURL, err := url.Parse(proxyAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
		}
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
		opts = append(opts, rpc.WithHTTPClient(httpClient))
	}

	rpcClient, err := rpc.DialOptions(context.Background(), rpcURL, opts...)
	if err != nil {
		return nil, err
	}

	return ethclient.NewClient(rpcClient), nil
}

// CheckERC20Balance 查询指定链上 ERC-20 余额是否 > 0
func (s *EthereumServiceImpl) CheckERC20Balance(chainID, walletAddress, contractAddress string) (bool, error) {
	s.mu.RLock()
	cc, ok := s.clients[chainID]
	s.mu.RUnlock()
	if !ok {
		return false, fmt.Errorf("unsupported chain: %s", chainID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()

	wallet := common.HexToAddress(walletAddress)
	contract := common.HexToAddress(contractAddress)

	// balanceOf(address) selector: 0x70a08231
	selector := []byte{0x70, 0xa0, 0x82, 0x31}
	paddedAddress := common.LeftPadBytes(wallet.Bytes(), 32)
	data := append(selector, paddedAddress...)

	msg := ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}

	result, err := cc.client.CallContract(ctx, msg, nil)
	if err != nil {
		return false, fmt.Errorf("on-chain query failed: %w", err)
	}

	balance := new(big.Int).SetBytes(result)
	return balance.Sign() > 0, nil
}

// GetChains 返回支持的链列表
func (s *EthereumServiceImpl) GetChains() []config.ChainConfig {
	return s.chains
}

// CheckERC721Ownership 查询指定链上 ERC-721 NFT 持有情况
func (s *EthereumServiceImpl) CheckERC721Ownership(chainID, walletAddress, contractAddress string) (bool, []*big.Int, error) {
	s.mu.RLock()
	cc, ok := s.clients[chainID]
	s.mu.RUnlock()
	if !ok {
		return false, nil, fmt.Errorf("unsupported chain: %s", chainID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()

	wallet := common.HexToAddress(walletAddress)
	contract := common.HexToAddress(contractAddress)

	// ERC-721 balanceOf(address) selector: 0x70a08231
	selector := []byte{0x70, 0xa0, 0x82, 0x31}
	paddedAddress := common.LeftPadBytes(wallet.Bytes(), 32)
	data := append(selector, paddedAddress...)

	result, err := cc.client.CallContract(ctx, ethereum.CallMsg{To: &contract, Data: data}, nil)
	if err != nil {
		return false, nil, fmt.Errorf("ERC-721 balanceOf query failed: %w", err)
	}

	balance := new(big.Int).SetBytes(result)
	if balance.Sign() == 0 {
		return false, nil, nil
	}

	// 安全检查：balanceOf 返回值过大说明可能不是 ERC-721 合约（如 ERC-20 余额）
	maxEnumerate := big.NewInt(1000)
	if balance.Cmp(maxEnumerate) > 0 {
		return true, nil, nil
	}

	// tokenOfOwnerByIndex(address, uint256) selector: 0x2f745c59
	count := int(balance.Int64())
	tokenIDs := make([]*big.Int, 0, count)
	indexSelector := []byte{0x2f, 0x74, 0x5c, 0x59}

	for i := 0; i < count; i++ {
		idx := common.LeftPadBytes(big.NewInt(int64(i)).Bytes(), 32)
		callData := append(indexSelector, paddedAddress...)
		callData = append(callData, idx...)

		tokenResult, err := cc.client.CallContract(ctx, ethereum.CallMsg{To: &contract, Data: callData}, nil)
		if err != nil {
			// 合约不支持 ERC721Enumerable（tokenOfOwnerByIndex revert），返回持有状态但不列出 tokenId
			return true, nil, nil
		}

		tokenID := new(big.Int).SetBytes(tokenResult)
		tokenIDs = append(tokenIDs, tokenID)
	}

	return true, tokenIDs, nil
}

// CheckERC1155Balance 查询指定链上 ERC-1155 代币持有数量
func (s *EthereumServiceImpl) CheckERC1155Balance(chainID, walletAddress, contractAddress, tokenID string) (bool, *big.Int, error) {
	s.mu.RLock()
	cc, ok := s.clients[chainID]
	s.mu.RUnlock()
	if !ok {
		return false, nil, fmt.Errorf("unsupported chain: %s", chainID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()

	wallet := common.HexToAddress(walletAddress)
	contract := common.HexToAddress(contractAddress)

	tid, ok := new(big.Int).SetString(tokenID, 10)
	if !ok {
		return false, nil, fmt.Errorf("invalid token ID: %s", tokenID)
	}

	// balanceOf(address,uint256) selector: 0x00fdd58e
	selector := []byte{0x00, 0xfd, 0xd5, 0x8e}
	paddedAddress := common.LeftPadBytes(wallet.Bytes(), 32)
	paddedTokenID := common.LeftPadBytes(tid.Bytes(), 32)
	data := append(selector, paddedAddress...)
	data = append(data, paddedTokenID...)

	result, err := cc.client.CallContract(ctx, ethereum.CallMsg{To: &contract, Data: data}, nil)
	if err != nil {
		return false, nil, fmt.Errorf("on-chain query failed: %w", err)
	}

	balance := new(big.Int).SetBytes(result)
	return balance.Sign() > 0, balance, nil
}

// Close 关闭所有连接
func (s *EthereumServiceImpl) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, cc := range s.clients {
		cc.client.Close()
	}
}
