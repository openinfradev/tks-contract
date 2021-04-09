package main

import (
	"github.com/openinfradev/tks-contract/pkg/contract"
	"github.com/openinfradev/tks-contract/pkg/log"
	"github.com/openinfradev/tks-proto/pbgo"
)

// InsertMockupContracts create mockup contract data in-memory.
func InsertMockupContracts(contract *contract.Accessor) error {
	mID, err := contractAccessor.Post("mock1", "tks-1000001", []string{"lma"},
		&pbgo.ContractQuota{
			Cpu:    14,
			Memory: 14,
		})
	if err != nil {
		return err
	}
	log.Info("Create Mockup data! new contract. mID:", mID)
	return nil
}
