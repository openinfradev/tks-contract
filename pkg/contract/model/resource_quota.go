package model

import (
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

// ResourceQuota represents a resource quota
type ResourceQuota struct {
	ID         uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Cpu        int64
	Memory     int64
	Block      int64
	BlockSsd   int64
	Fs         int64
	FsSsd      int64
	ContractID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (r *ResourceQuota) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return nil
}
