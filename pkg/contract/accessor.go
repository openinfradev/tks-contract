package contract

import (
	"fmt"
	"time"

	pb "github.com/openinfradev/tks-proto/pbgo"
)

// Accessor is an accessor to in-memory contracts.
type Accessor struct {
	contracts map[ID]Contract
}

// NewContractAccessor returns new contract accessor's ptr.
func NewContractAccessor() *Accessor {
	return &Accessor{
		contracts: map[ID]Contract{},
	}
}

// Get returns contract data from in-memory map.
func (c Accessor) Get(id ID) (*Contract, error) {
	contract, exists := c.contracts[id]
	if !exists {
		return &Contract{}, fmt.Errorf("Not found contract for %s", id)
	}

	return &contract, nil
}

// Post inserts new contract in in-memory map if contractId does not exist.
func (c *Accessor) Post(contractorName string, id ID,
	availableServices []string, quota *pb.ContractQuota) (McOpsID, error) {
	if _, exists := c.contracts[id]; exists {
		return McOpsID{}, fmt.Errorf("Already exists contractId %s", id)
	}
	newMcOpsID := GenerateMcOpsID()
	c.contracts[id] = Contract{
		ContractorName:    contractorName,
		ID:                id,
		AvailableServices: availableServices,
		Quota:             quota,
		LastUpdatedTs:     &LastUpdatedTime{time.Now()},
		McOpsID:           newMcOpsID,
	}
	return newMcOpsID, nil
}

// Update updates contract data by contractID.
// Not implemented yet.
func (c *Accessor) Update(contractorName string, id ID,
	availableServices []string, quota *pb.ContractQuota) error {
	if _, exists := c.contracts[id]; exists {
		return fmt.Errorf("Already exists contractId %s", id)
	}
	c.contracts[id] = Contract{
		ContractorName:    contractorName,
		ID:                id,
		AvailableServices: availableServices,
		Quota:             quota,
		LastUpdatedTs:     &LastUpdatedTime{time.Now()},
	}
	return nil
}
