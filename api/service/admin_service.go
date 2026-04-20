package service

import (
	"fmt"

	"chain-access/api/model"
	"chain-access/api/repository"
)

// AdminService admin business logic interface
type AdminService interface {
	VerifyAdminAccess(address, signature, chainID, nftContract string) (string, error)
	ListBalances(page, size int, address string) (*model.BalanceListResponse, error)
}

type adminServiceImpl struct {
	authService AuthService
	ethService  EthereumService
	balanceRepo repository.BalanceRepository
}

// NewAdminService creates admin service
func NewAdminService(
	authSvc AuthService,
	ethSvc EthereumService,
	balanceRepo repository.BalanceRepository,
) AdminService {
	return &adminServiceImpl{
		authService: authSvc,
		ethService:  ethSvc,
		balanceRepo: balanceRepo,
	}
}

func (s *adminServiceImpl) VerifyAdminAccess(address, signature, chainID, nftContract string) (string, error) {
	token, err := s.authService.VerifySignature(address, signature)
	if err != nil {
		return "", fmt.Errorf("signature verification failed: %w", err)
	}

	hasNFT, _, err := s.ethService.CheckERC721Ownership(chainID, address, nftContract)
	if err != nil {
		return "", fmt.Errorf("NFT check failed: %w", err)
	}
	if !hasNFT {
		return "", fmt.Errorf("admin access denied: no NFT found on chain %s contract %s", chainID, nftContract)
	}

	return token, nil
}

func (s *adminServiceImpl) ListBalances(page, size int, address string) (*model.BalanceListResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	data, total, err := s.balanceRepo.List(page, size, address)
	if err != nil {
		return nil, err
	}
	return &model.BalanceListResponse{
		Data:  data,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

// noopAdminService is used when MySQL is not configured
type noopAdminService struct{}

func NewNoopAdminService() AdminService { return &noopAdminService{} }

func (s *noopAdminService) VerifyAdminAccess(address, signature, chainID, nftContract string) (string, error) {
	return "", fmt.Errorf("admin features not configured")
}

func (s *noopAdminService) ListBalances(page, size int, address string) (*model.BalanceListResponse, error) {
	return &model.BalanceListResponse{Data: []model.UserBalance{}, Total: 0, Page: page, Size: size}, nil
}
