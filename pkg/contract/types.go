package contract

import (
	"time"

	uuid "github.com/google/uuid"
	"github.com/lib/pq"
)

// Contract represents a contract data
// type Contract struct {
// 	ID                uuid.UUID `gorm:"primarykey"`
// 	ContractorName    string
// 	AvailableServices []string
// 	CspID             uuid.UUID
// 	UpdatedAt         time.Time
// 	CreatedAt         time.Time
// }

// Contract represents a contract data in Database.
type Contract struct {
	ID                uuid.UUID
	ContractorName    string
	AvailableServices pq.StringArray
	CspID             uuid.UUID
	UpdatedAt         time.Time
	CreatedAt         time.Time
	Quota             *ResourceQuota
}

// ResourceQuota represents a resource quota
type ResourceQuota struct {
	ID                uuid.UUID
	UpdatedAt         time.Time
	CreatedAt         time.Time
	ContractorName    string
	AvailableServices []string
	CspID             uuid.UUID
}
