package contract

import (
	uuid "github.com/google/uuid"
)

// ResourceQuotaParam is parameter to update quota for contract.
type ResourceQuotaParam struct {
	ID         uuid.UUID
	Cpu        int64
	Memory     int64
	Block      int64
	BlockSsd   int64
	Fs         int64
	FsSsd      int64
	ContractID uuid.UUID
}
