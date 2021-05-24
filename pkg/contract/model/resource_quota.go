package model

import (
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

// ResourceQuota represents a resource quota
type ResourceQuota struct {
	ID                uuid.UUID `gorm:"primarykey"`
	UpdatedAt         time.Time
	CreatedAt         time.Time
	ContractorName    string
	AvailableServices []string
	CspID             uuid.UUID
}

func (r *ResourceQuota) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return nil
}
