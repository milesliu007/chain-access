package service

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"chain-access/api/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthService 认证服务接口
type AuthService interface {
	// GenerateChallenge 为指定地址生成 challenge
	GenerateChallenge(address string) string
	// VerifySignature 验证签名并返回 JWT
	VerifySignature(address, signature string) (string, error)
	// ValidateJWT 验证 JWT 并返回地址
	ValidateJWT(tokenString string) (string, error)
}

// authServiceImpl 认证服务实现
type authServiceImpl struct {
	challengeRepo repository.ChallengeRepository
	jwtSecret     []byte
}

// NewAuthService 创建认证服务
func NewAuthService(jwtSecret []byte, challengeRepo repository.ChallengeRepository) AuthService {
	return &authServiceImpl{
		jwtSecret:     jwtSecret,
		challengeRepo: challengeRepo,
	}
}

// GenerateChallenge 为指定地址生成 challenge
func (s *authServiceImpl) GenerateChallenge(address string) string {
	nonce := uuid.New().String()
	message := fmt.Sprintf("Sign this message to authenticate with chain-access: %s", nonce)

	addr := strings.ToLower(address)
	s.challengeRepo.Store(addr, repository.ChallengeEntry{
		Message:   message,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	return message
}

// VerifySignature 验证签名并返回 JWT
func (s *authServiceImpl) VerifySignature(address, signature string) (string, error) {
	addr := strings.ToLower(address)

	// 取出 challenge（一次性使用）
	entry, ok := s.challengeRepo.LoadAndDelete(addr)
	if !ok {
		return "", fmt.Errorf("challenge not found for this address, please request /auth/challenge first")
	}

	// 检查 TTL
	if time.Now().After(entry.ExpiresAt) {
		return "", fmt.Errorf("challenge expired, please request a new one")
	}

	// EIP-191 签名验证
	recoveredAddr, err := recoverAddress(entry.Message, signature)
	if err != nil {
		return "", fmt.Errorf("signature verification failed: %w", err)
	}

	if strings.ToLower(recoveredAddr.Hex()) != addr {
		return "", fmt.Errorf("signature address mismatch")
	}

	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"address": addr,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT 验证 JWT 并返回地址
func (s *authServiceImpl) ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unsupported signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid JWT")
	}

	address, ok := claims["address"].(string)
	if !ok {
		return "", fmt.Errorf("missing address in JWT")
	}

	return address, nil
}

// recoverAddress 从 EIP-191 签名中恢复地址
func recoverAddress(message, sig string) (common.Address, error) {
	sigBytes, err := hex.DecodeString(strings.TrimPrefix(sig, "0x"))
	if err != nil {
		return common.Address{}, fmt.Errorf("invalid signature format: %w", err)
	}

	if len(sigBytes) != 65 {
		return common.Address{}, fmt.Errorf("invalid signature length: expected 65 bytes, got %d", len(sigBytes))
	}

	// MetaMask 签名 v=27/28，go-ethereum 期望 0/1
	v := sigBytes[64]
	if v == 27 || v == 28 {
		sigBytes[64] = v - 27
	} else if v != 0 && v != 1 {
		return common.Address{}, fmt.Errorf("invalid signature v value: %d, expected 0/1/27/28", v)
	}

	// EIP-191 前缀
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))

	pubKey, err := crypto.SigToPub(hash.Bytes(), sigBytes)
	if err != nil {
		return common.Address{}, fmt.Errorf("ecrecover failed: %w", err)
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}

// IsValidAddress 校验以太坊地址格式
func IsValidAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") && !strings.HasPrefix(address, "0X") {
		return false
	}
	addr := strings.TrimPrefix(strings.TrimPrefix(address, "0x"), "0X")
	if len(addr) != 40 {
		return false
	}
	_, err := hex.DecodeString(addr)
	return err == nil
}
