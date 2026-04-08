package repository

import (
	"sync"
	"time"
)

const challengeTTL = 5 * time.Minute

// ChallengeEntry challenge 条目
type ChallengeEntry struct {
	Message   string
	ExpiresAt time.Time
}

// ChallengeRepository challenge 存储接口
type ChallengeRepository interface {
	// Store 存储 challenge
	Store(address string, entry ChallengeEntry)
	// LoadAndDelete 取出并删除 challenge（一次性使用）
	LoadAndDelete(address string) (ChallengeEntry, bool)
}

// MemoryChallengeRepository 基于内存的 challenge 存储
type MemoryChallengeRepository struct {
	store sync.Map
}

// NewMemoryChallengeRepository 创建内存 challenge 存储，启动后台清理
func NewMemoryChallengeRepository() *MemoryChallengeRepository {
	r := &MemoryChallengeRepository{}
	go r.cleanupExpired()
	return r
}

// Store 存储 challenge
func (r *MemoryChallengeRepository) Store(address string, entry ChallengeEntry) {
	r.store.Store(address, entry)
}

// LoadAndDelete 取出并删除 challenge
func (r *MemoryChallengeRepository) LoadAndDelete(address string) (ChallengeEntry, bool) {
	val, ok := r.store.LoadAndDelete(address)
	if !ok {
		return ChallengeEntry{}, false
	}
	return val.(ChallengeEntry), true
}

// cleanupExpired 定期清理过期 challenge，防止内存泄漏
func (r *MemoryChallengeRepository) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		r.store.Range(func(key, value any) bool {
			entry := value.(ChallengeEntry)
			if now.After(entry.ExpiresAt) {
				r.store.Delete(key)
			}
			return true
		})
	}
}
