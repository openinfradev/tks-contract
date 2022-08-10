package contract

import (
	"fmt"

	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-common/pkg/log"
	model "github.com/openinfradev/tks-contract/pkg/contract/model"
	pb "github.com/openinfradev/tks-proto/tks_pb"
)

// Accessor is an accessor to in-memory contracts.
type Accessor struct {
	db *gorm.DB
}

// New returns new accessor's ptr.
func New(db *gorm.DB) *Accessor {
	return &Accessor{
		db: db,
	}
}

// GetContract returns a contract from database.
func (x *Accessor) GetContract(id string) (*pb.Contract, error) {
	var contract model.Contract
	res := x.db.First(&contract, "id = ?", id)
	if res.RowsAffected == 0 || res.Error != nil {
		return &pb.Contract{}, fmt.Errorf("Not found contract for %s", id)
	}
	quota, err := x.GetResourceQuota(contract.ID)
	if err != nil {
		return &pb.Contract{}, err
	}
	resContract := reflectToPbContract(contract, &quota)
	return &resContract, nil
}

// GetDefaultContract returns a contract from database.
func (x *Accessor) GetDefaultContract() (*pb.Contract, error) {
	var contract model.Contract
	res := x.db.First(&contract, "contractor_name = 'default'")
	if res.RowsAffected == 0 || res.Error != nil {
		return &pb.Contract{}, fmt.Errorf("Not found default contract")
	}
	quota, err := x.GetResourceQuota(contract.ID)
	if err != nil {
		return &pb.Contract{}, err
	}
	resContract := reflectToPbContract(contract, &quota)
	return &resContract, nil
}

// getContract returns a resource quota from database.
func (x *Accessor) GetResourceQuota(contractID string) (pb.ContractQuota, error) {
	var quota model.ResourceQuota
	res := x.db.Limit(1).Find(&quota, "contract_id = ?", contractID)
	if res.RowsAffected == 0 || res.Error != nil {
		return pb.ContractQuota{}, fmt.Errorf("Not found quota for contract id %s", contractID)
	}

	return reflectToPbQuota(quota), nil
}

// List returns a list of contracts from database.
func (x *Accessor) List(offset, limit int) ([]*pb.Contract, error) {
	var (
		contracts       []model.Contract
		quota           model.ResourceQuota
		resultContracts []*pb.Contract
	)
	res := x.db.Offset(offset).Limit(limit).Find(&contracts)
	if res.Error != nil {
		return nil, res.Error
	}
	for _, contract := range contracts {
		quota = model.ResourceQuota{}
		res = x.db.Limit(1).Find(&quota, "contract_id = ?", contract.ID)
		if res.RowsAffected == 0 || res.Error != nil {
			return nil, fmt.Errorf("Not found quota for contract id %s", contract.ID)
		}
		pbQuota := reflectToPbQuota(quota)
		resContract := reflectToPbContract(contract, &pbQuota)
		resultContracts = append(resultContracts, &resContract)
	}
	return resultContracts, nil
}

// Create creates a new contract in database.
func (x *Accessor) Create(name string, availableServices []string, quota *pb.ContractQuota) (string, error) {
	pqStrArr := pq.StringArray{}

	for _, svc := range availableServices {
		pqStrArr = append(pqStrArr, svc)
	}

	contract := model.Contract{ContractorName: name, AvailableServices: pqStrArr}
	err := x.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Create(&contract)
		if res.Error != nil {
			return res.Error
		}
		res = tx.Create(&model.ResourceQuota{Cpu: quota.Cpu, Memory: quota.Memory,
			Block: quota.Block, BlockSsd: quota.BlockSsd, Fs: quota.Fs, FsSsd: quota.FsSsd, ContractID: contract.ID})
		if res.Error != nil {
			return res.Error
		}
		log.Info("sucessfully created contract ID ", contract.ID)
		return nil
	})

	return contract.ID, err
}

// Delete contract
func (x *Accessor) Delete(contractId string) error {
	err := x.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Delete(&model.ResourceQuota{}, "contract_id = ?", contractId)
		log.Info("resource quota is deleted! contractId : ", contractId)
		if res.Error != nil {
			return fmt.Errorf("could not delete resource quota for contractId %s", contractId)
		}

		res = tx.Delete(&model.Contract{}, "id = ?", contractId)
		log.Info("contract is deleted! contractId : ", contractId)
		if res.Error != nil {
			return fmt.Errorf("could not delete contract for contractId %s", contractId)
		}
		return nil
	})

	return err
}

// UpdateResourceQuota updates resource quota.
func (x *Accessor) UpdateResourceQuota(contractID string, quota *pb.ContractQuota) (
	p *pb.ContractQuota, c *pb.ContractQuota, err error) {
	prev, err := x.GetResourceQuota(contractID)
	if err != nil {
		return &pb.ContractQuota{}, &pb.ContractQuota{}, fmt.Errorf("not found resource quota for contract ID %s", contractID)
	}

	values := map[string]interface{}{
		"cpu":       prev.Cpu,
		"memory":    prev.Memory,
		"block":     prev.Block,
		"block_ssd": prev.BlockSsd,
		"fs":        prev.Fs,
		"fs_ssd":    prev.FsSsd,
	}

	if quota.Cpu != 0 {
		values["cpu"] = quota.Cpu
	}
	if quota.Memory != 0 {
		values["memory"] = quota.Memory
	}
	if quota.Block != 0 {
		values["block"] = quota.Block
	}
	if quota.Block != 0 {
		values["block_ssd"] = quota.BlockSsd
	}
	if quota.Fs != 0 {
		values["fs"] = quota.Fs
	}
	if quota.FsSsd != 0 {
		values["fs_ssd"] = quota.FsSsd
	}

	res := x.db.Model(&model.ResourceQuota{}).
		Where("contract_id = ?", contractID).
		Updates(values)

	if res.Error != nil || res.RowsAffected == 0 {
		return nil, nil, fmt.Errorf("nothing updated in resource_quota for contract id %s", contractID)
	}

	curr, err := x.GetResourceQuota(contractID)
	return &prev, &curr, err
}

// UpdateAvailableServices updates available service list and resource quota.
func (x *Accessor) UpdateAvailableServices(id string, availableServices []string) (
	prev []string, curr []string, err error) {
	pqStrArr := pq.StringArray{}

	for _, svc := range availableServices {
		pqStrArr = append(pqStrArr, svc)
	}
	var (
		contract model.Contract
	)
	if res := x.db.First(&contract, "id = ?", id); res.RowsAffected == 0 || res.Error != nil {
		return nil, nil, fmt.Errorf("could not find contract for contract id %s", id)
	}
	prev = contract.AvailableServices
	if res := x.db.Model(&model.Contract{}).Where("id = ?", id).Update("available_services", pqStrArr); res.RowsAffected == 0 || res.Error != nil {
		return prev, curr, fmt.Errorf("RowsAffected is 0 for contract id %s", id)
	}

	if res := x.db.First(&contract, "id = ?", id); res.RowsAffected == 0 || res.Error != nil {
		return nil, nil, fmt.Errorf("could not find contract for contract id %s", id)
	}
	curr = contract.AvailableServices
	return prev, curr, nil
}

func reflectToPbContract(contract model.Contract, quota *pb.ContractQuota) pb.Contract {
	return pb.Contract{
		ContractId:        contract.ID,
		ContractorName:    contract.ContractorName,
		AvailableServices: contract.AvailableServices,
		UpdatedAt:         timestamppb.New(contract.UpdatedAt),
		CreatedAt:         timestamppb.New(contract.CreatedAt),
		Quota:             quota,
	}
}

func reflectToPbQuota(quota model.ResourceQuota) pb.ContractQuota {
	return pb.ContractQuota{
		Cpu:      quota.Cpu,
		Memory:   quota.Memory,
		Block:    quota.Block,
		BlockSsd: quota.BlockSsd,
		Fs:       quota.Fs,
		FsSsd:    quota.FsSsd,
	}
}
