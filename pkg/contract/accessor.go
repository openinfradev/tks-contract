package contract

import (
	"fmt"

	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"

	uuid "github.com/google/uuid"
	model "github.com/sktelecom/tks-contract/pkg/contract/model"
	pb "github.com/sktelecom/tks-proto/pbgo"
	"gorm.io/gorm"
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
func (x *Accessor) GetContract(id uuid.UUID) (*pb.Contract, error) {
	var contract model.Contract
	res := x.db.First(&contract, id)
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

// getContract returns a resource quota from database.
func (x *Accessor) GetResourceQuota(contractID uuid.UUID) (pb.ContractQuota, error) {
	var quota model.ResourceQuota
	res := x.db.First(&quota, "contract_id = ?", contractID)
	if res.RowsAffected == 0 || res.Error != nil {
		return pb.ContractQuota{}, fmt.Errorf("Not found quota for contract id %s", contractID)
	}

	return reflectToPbQuota(quota), nil
}

// List returns a list of contracts from database.
func (x *Accessor) List(offset, limit int) ([]pb.Contract, error) {
	var (
		contracts       []model.Contract
		quota           model.ResourceQuota
		resultContracts []pb.Contract
	)
	res := x.db.Offset(offset).Limit(limit).Find(&contracts)
	if res.Error != nil {
		return nil, res.Error
	}
	for _, contract := range contracts {
		res = x.db.First(&quota, "contract_id = ?", contract.ID)
		if res.RowsAffected == 0 || res.Error != nil {
			return nil, fmt.Errorf("Not found quota for contract id %s", contract.ID)
		}
		pbQuota := reflectToPbQuota(quota)
		resultContracts = append(resultContracts, reflectToPbContract(contract, &pbQuota))
	}
	return resultContracts, nil
}

// Create creates a new contract in database.
func (x *Accessor) Create(name string, availableServices []string, quota ResourceQuotaParam) (uuid.UUID, error) {
	pqStrArr := pq.StringArray{}

	for _, svc := range availableServices {
		pqStrArr = append(pqStrArr, svc)
	}

	var contract model.Contract
	err := x.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Create(&model.Contract{ContractorName: name, AvailableServices: pqStrArr})
		if res.Error != nil {
			return res.Error
		}
		res = tx.First(&contract, "contractor_name = ?", name)
		if res.Error != nil {
			return res.Error
		}
		res = tx.Create(&model.ResourceQuota{Cpu: quota.Cpu, Memory: quota.Memory,
			Block: quota.Block, BlockSsd: quota.BlockSsd, Fs: quota.Fs, FsSsd: quota.FsSsd, ContractID: contract.ID})
		if res.Error != nil {
			return res.Error
		}
		return nil
	})

	return contract.ID, err
}

// UpdateResourceQuota updates resource quota.
func (x *Accessor) UpdateResourceQuota(contractID uuid.UUID, quota ResourceQuotaParam) (
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

	curr, err := x.GetResourceQuota(contractID)

	return &prev, &curr, res.Error
}

// UpdateAvailableServices updates available service list and resource quota.
func (x *Accessor) UpdateAvailableServices(id uuid.UUID, availableServices []string) (
	prev []string, curr []string, err error) {
	pqStrArr := pq.StringArray{}

	for _, svc := range availableServices {
		pqStrArr = append(pqStrArr, svc)
	}
	var (
		contract model.Contract
	)
	if res := x.db.First(&contract, id); res.RowsAffected == 0 || res.Error != nil {
		return nil, nil, fmt.Errorf("not exist contract for contract id %s", id)
	}
	prev = contract.AvailableServices
	if res := x.db.Model(&model.Contract{}).Where("id = ?", id).Update("available_services", pqStrArr); res.RowsAffected == 0 || res.Error != nil {
		return prev, curr, fmt.Errorf("RowsAffected is 0 for contract id %s", id)
	}

	if res := x.db.First(&contract, id); res.RowsAffected == 0 || res.Error != nil {
		return nil, nil, fmt.Errorf("not exist contract for contract id %s", id)
	}
	curr = contract.AvailableServices
	return prev, curr, nil
}

func reflectToPbContract(contract model.Contract, quota *pb.ContractQuota) pb.Contract {
	return pb.Contract{
		ContractId:        contract.ID.String(),
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
