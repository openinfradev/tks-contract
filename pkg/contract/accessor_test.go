package contract_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/openinfradev/tks-contract/pkg/contract"
	pb "github.com/openinfradev/tks-proto/pbgo"
)

var contractID uuid.UUID

func getAccessor() (*contract.Accessor, error) {
	dsn := "host=localhost user=postgres password=password dbname=tks port=5432 sslmode=disable TimeZone=Asia/Seoul"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return contract.New(db), nil
}
func TestCreateContract(t *testing.T) {
	accessor, err := getAccessor()
	if err != nil {
		t.Errorf("an error was unexpected while initilizing database %s", err)
	}
	quota := pb.ContractQuota{
		Cpu:    256,
		Memory: 12800000,
		Block:  12800000,
		Fs:     12800000,
	}
	contractName := getRandomString("gotest")
	contractID, err = accessor.Create(contractName, []string{"lma"}, &quota)
	if err != nil {
		t.Errorf("an error was unexpected while creating new contract: %s", err)
	}
	t.Logf("new contract id: %s", contractID)
}

func TestUpdateAvailableServices(t *testing.T) {
	accessor, err := getAccessor()
	if err != nil {
		t.Errorf("an error was unexpected while initilizing database %s", err)
	}
	_, _, err = accessor.UpdateAvailableServices(contractID, []string{"lma", "sm"})
	if err != nil {
		t.Errorf("an error was unexpected while querying contract data %s", err)
	}
}

func TestUpdateResourceQuota(t *testing.T) {
	accessor, err := getAccessor()
	if err != nil {
		t.Errorf("an error was unexpected while initilizing database %s", err)
	}
	quota := pb.ContractQuota{
		Cpu:    128,
		Memory: 1280000,
	}
	_, _, err = accessor.UpdateResourceQuota(contractID, &quota)

	if err != nil {
		t.Errorf("an error was unexpected while querying contract data %s", err)
	}
}
func TestGetContract(t *testing.T) {
	accessor, err := getAccessor()
	if err != nil {
		t.Errorf("an error was unexpected while initilizing database %s", err)
	}
	contract, err := accessor.GetContract(contractID)

	if err != nil {
		t.Errorf("an error was unexpected while querying contract data %s", err)
	}

	t.Logf("contractor name: %s", contract.ContractorName)
	t.Logf("quota cpu: %d", contract.Quota.Cpu)
}

func getRandomString(prefix string) string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return fmt.Sprintf("%s-%d", prefix, r.Int31n(1000000000))
}
