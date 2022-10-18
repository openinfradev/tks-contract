package model

import (
	"time"

	uuid "github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-common/pkg/helper"
)

// Contract represents a contract data in Database.
type Contract struct {
	ID                string         `gorm:"primaryKey"`
	ContractorName    string         `gorm:"uniqueIndex"`
	AvailableServices pq.StringArray `gorm:"type:text[]"`
	Creator           uuid.UUID
	Description       string
	UpdatedAt         time.Time
	CreatedAt         time.Time
}

func (c *Contract) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = helper.GenerateContractId()
	return nil
}
