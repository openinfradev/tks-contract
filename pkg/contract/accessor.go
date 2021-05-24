package contract

import (
	"fmt"

	"github.com/lib/pq"

	uuid "github.com/google/uuid"
	model "github.com/sktelecom/tks-contract/pkg/contract/model"
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

// Get returns a contract from database.
func (x *Accessor) Get(id uuid.UUID) (model.Contract, error) {
	var contract model.Contract
	result := x.db.First(&contract, id)
	if result.RowsAffected == 0 || result.Error != nil {
		return model.Contract{}, fmt.Errorf("Not found contract for %s", id)
	}

	return contract, nil
}

// List returns a list of contracts from database.
func (x *Accessor) List(offset, limit int) ([]model.Contract, error) {
	var contracts []model.Contract
	result := x.db.Offset(offset).Limit(limit).Find(&contracts)
	return contracts, result.Error
}

// Create creates a new contract in database.
func (x *Accessor) Create(name string, availableServices []string, cspID uuid.UUID) (uuid.UUID, error) {
	pqStrArr := pq.StringArray{}

	for _, svc := range availableServices {
		pqStrArr = append(pqStrArr, svc)
	}

	var contract model.Contract
	err := x.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Create(&model.Contract{ContractorName: name, AvailableServices: pqStrArr, CspID: cspID})
		if res.Error != nil {
			return res.Error
		}
		res = tx.First(&contract, "contractor_name = ?", name)
		if res.Error != nil {
			return res.Error
		}
		return nil
	})

	return contract.ID, err
}

// Update updates available service list and resource quota.
func (x *Accessor) Update(id uuid.UUID, availableServices []string) error {
	pqStrArr := pq.StringArray{}

	for _, svc := range availableServices {
		pqStrArr = append(pqStrArr, svc)
	}

	err := x.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&model.Contract{}).
			Where("id = ?", id).
			Update("available_services", pqStrArr)
		if res.Error != nil {
			return res.Error
		}
		return nil
	})

	return err
}
