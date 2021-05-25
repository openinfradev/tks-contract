package model

import (
	"time"

	uuid "github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Contract represents a contract data in Database.
type Contract struct {
	ID                uuid.UUID      `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	ContractorName    string         `gorm:"uniqueIndex"`
	AvailableServices pq.StringArray `gorm:"type:text[]"`
	UpdatedAt         time.Time
	CreatedAt         time.Time
}

func (c *Contract) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return nil
}
