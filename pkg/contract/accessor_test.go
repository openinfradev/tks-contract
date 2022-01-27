package contract_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	log "github.com/openinfradev/tks-common/pkg/log"
	"github.com/openinfradev/tks-contract/pkg/contract"
	model "github.com/openinfradev/tks-contract/pkg/contract/model"
	pb "github.com/openinfradev/tks-proto/tks_pb"

	helper "github.com/openinfradev/tks-common/pkg/helper"
)

var (
	err error
	testDBHost string
	testDBPort string
)

var (
	contractId uuid.UUID
)

func init() {
	log.Disable()
}

func getAccessor() (*contract.Accessor, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		testDBHost, "postgres", "password", "tks", testDBPort )
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`);

	if err := db.AutoMigrate(&model.Contract{}); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.ResourceQuota{}); err != nil {
		return nil, err
	}

	return contract.New(db), nil
}

func TestMain(m *testing.M) {
	pool, resource, err := helper.CreatePostgres()
	if err != nil {
		fmt.Printf("Could not create postgres: %s", err)
		os.Exit(-1)
	}
	testDBHost, testDBPort = helper.GetHostAndPort(resource)

	code := m.Run()

	if err := helper.RemovePostgres(pool, resource); err != nil {
		fmt.Printf("Could not remove postgres: %s", err)
		os.Exit(-1)
	}
	os.Exit(code)
}


// TestCases

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
	contractId, err = accessor.Create(contractName, []string{"lma"}, &quota)
	if err != nil {
		t.Errorf("an error was unexpected while creating new contract: %s", err)
	}
	t.Logf("new contract id: %s", contractId)
}

func TestUpdateAvailableServices(t *testing.T) {
	accessor, err := getAccessor()
	if err != nil {
		t.Errorf("an error was unexpected while initilizing database %s", err)
	}
	_, _, err = accessor.UpdateAvailableServices(contractId, []string{"lma", "sm"})
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
	_, _, err = accessor.UpdateResourceQuota(contractId, &quota)

	if err != nil {
		t.Errorf("an error was unexpected while querying contract data %s", err)
	}
}
func TestGetContract(t *testing.T) {
	accessor, err := getAccessor()
	if err != nil {
		t.Errorf("an error was unexpected while initilizing database %s", err)
	}
	contract, err := accessor.GetContract(contractId)

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
