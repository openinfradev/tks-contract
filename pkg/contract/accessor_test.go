package contract_test

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sktelecom/tks-contract/pkg/contract"
)

func getAccessor() (*contract.Accessor, error) {
	dsn := "host=localhost user=postgres password=password dbname=tks port=5432 sslmode=disable TimeZone=Asia/Seoul"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return contract.New(db), nil
}
func TestCreate(t *testing.T) {
	accessor, err := getAccessor()
	if err != nil {
		t.Errorf("an error was unexpected while initilizing database %s", err)
	}
	contractID, err := accessor.Create("tester2", []string{"lma"}, uuid.MustParse("677dee5f-f224-482e-8aad-fd312cba19fe"))
	if err != nil {
		t.Errorf("an error was unexpected while creating new contract: %s", err)
	}
	t.Logf("new contract id: %s", contractID)
}

func TestUpdate(t *testing.T) {
	accessor, err := getAccessor()
	if err != nil {
		t.Errorf("an error was unexpected while initilizing database %s", err)
	}
	err = accessor.Update(uuid.MustParse("edcaa975-dde4-4c4d-94f7-36bc38fe7064"),
		[]string{"lma", "sm"})

	if err != nil {
		t.Errorf("an error was unexpected while querying contract data %s", err)
	}
}

func TestGet(t *testing.T) {
	accessor, err := getAccessor()
	if err != nil {
		t.Errorf("an error was unexpected while initilizing database %s", err)
	}
	contract, err := accessor.Get(uuid.MustParse("edcaa975-dde4-4c4d-94f7-36bc38fe7064"))

	if err != nil {
		t.Errorf("an error was unexpected while querying contract data %s", err)
	}

	t.Logf("contractor name: %s", contract.ContractorName)
}
