package repository

import (
	"fmt"

	"chain-access/api/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// BalanceRepository MySQL balance query interface
type BalanceRepository interface {
	List(page, size int, address string) ([]model.UserBalance, int64, error)
}

type balanceRepo struct {
	db *gorm.DB
}

// NewBalanceRepository creates a GORM-backed balance repository
func NewBalanceRepository(dsn string) (BalanceRepository, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect mysql: %w", err)
	}
	return &balanceRepo{db: db}, nil
}

func (r *balanceRepo) List(page, size int, address string) ([]model.UserBalance, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	q := r.db.Model(&model.UserBalance{})
	if address != "" {
		q = q.Where("address LIKE ?", "%"+address+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	var list []model.UserBalance
	if err := q.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	if list == nil {
		list = []model.UserBalance{}
	}
	return list, total, nil
}
