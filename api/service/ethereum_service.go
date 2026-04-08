package service

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const rpcTimeout = 10 * time.Second

// EthereumService 链上查询服务接口
type EthereumService interface {
	// CheckERC20Balance 查询 ERC-20 余额是否 > 0
	CheckERC20Balance(walletAddress, contractAddress string) (bool, error)
	// Close 关闭连接
	Close()
}

// ethereumServiceImpl 链上查询服务实现
type ethereumServiceImpl struct {
	client *ethclient.Client
}

// NewEthereumService 创建链上查询服务
func NewEthereumService(infuraURL string) (EthereumService, error) {
	proxyURL, err := url.Parse("http://127.0.0.1:7897")
	if err != nil {
		return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	rpcClient, err := rpc.DialOptions(context.Background(), infuraURL, rpc.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Infura: %w", err)
	}

	client := ethclient.NewClient(rpcClient)
	return &ethereumServiceImpl{client: client}, nil
}

// CheckERC20Balance 查询 ERC-20 余额是否 > 0
func (s *ethereumServiceImpl) CheckERC20Balance(walletAddress, contractAddress string) (bool, error) {
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

	result, err := s.client.CallContract(ctx, msg, nil)
	if err != nil {
		return false, fmt.Errorf("on-chain query failed: %w", err)
	}

	balance := new(big.Int).SetBytes(result)
	return balance.Sign() > 0, nil
}

// Close 关闭连接
func (s *ethereumServiceImpl) Close() {
	s.client.Close()
}
