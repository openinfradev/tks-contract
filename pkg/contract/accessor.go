package contract

import (
	"errors"
	"fmt"
	"time"

	pb "github.com/openinfradev/tks-proto/pbgo"
)

type ContractAccessor struct {
	contracts map[ContractId]Contract
}

func NewContractAccessor() *ContractAccessor {
	return &ContractAccessor{
		contracts: map[ContractId]Contract{},
	}
}

func (c ContractAccessor) Get(contractId ContractId) (*Contract, error) {
	contract, exists := c.contracts[contractId]
	if !exists {
		return &Contract{}, errors.New(fmt.Sprintf("Not found contract for %s", contractId))
	}

	return &contract, nil
}

func (c *ContractAccessor) Put(contractorName string, contractId ContractId,
	availableServices []string, quota *pb.ContractQuota) (McOpsId, error) {
	if _, exists := c.contracts[contractId]; exists {
		return McOpsId{}, errors.New(fmt.Sprintf("Already exists contractId %s", contractId))
	}
	newMcOpsId := GenerateMcOpsId()
	c.contracts[contractId] = Contract{
		ContractorName:    contractorName,
		ContractId:        contractId,
		AvailableServices: availableServices,
		Quota:             quota,
		LastUpdatedTs:     &LastUpdatedTime{time.Now()},
		McOpsId:           newMcOpsId,
	}
	return newMcOpsId, nil
}

func (c *ContractAccessor) Update(contractorName string, contractId ContractId,
	availableServices []string, quota *pb.ContractQuota) error {
	if _, exists := c.contracts[contractId]; exists {
		return errors.New(fmt.Sprintf("Already exists contractId %s", contractId))
	}
	c.contracts[contractId] = Contract{
		ContractorName:    contractorName,
		ContractId:        contractId,
		AvailableServices: availableServices,
		Quota:             quota,
		LastUpdatedTs:     &LastUpdatedTime{time.Now()},
	}
	return nil
}
