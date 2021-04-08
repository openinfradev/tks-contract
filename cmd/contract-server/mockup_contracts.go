package main

import (
	"log"

	"github.com/openinfradev/tks-contract/pkg/contract"
	"github.com/openinfradev/tks-proto/pbgo"
)

// InsertMockupContracts create mockup contract data in-memory.
func InsertMockupContracts(contract *contract.ContractAccessor) error {
	mID, err := contractAccessor.Post("mock1", "tks-1000001", []string{"lma"},
		&pbgo.ContractQuota{
			Cpu:    14,
			Memory: 14,
		})
	if err != nil {
		return err
	}
	log.Println("Create new contract. mID:", mID)
	return nil
}
