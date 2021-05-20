package main

import (
	"github.com/sktelecom/tks-contract/pkg/contract"
	"github.com/sktelecom/tks-contract/pkg/log"
	"github.com/sktelecom/tks-proto/pbgo"
)

// InsertMockupContracts create mockup contract data in-memory.
func InsertMockupContracts(contract *contract.Accessor) error {
	err := contractAccessor.Post("mock1", "tks-1000001", "xxxx-xxxxxxx-xxxxx-xxxx",
		[]string{"lma"},
		&pbgo.ContractQuota{
			Cpu:    14,
			Memory: 14,
		})
	if err != nil {
		return err
	}
	log.Info("Create Mockup data!")
	return nil
}
