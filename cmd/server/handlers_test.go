package main

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sktelecom/tks-contract/pkg/contract"
	gc "github.com/sktelecom/tks-contract/pkg/grpc-client"
	pb "github.com/sktelecom/tks-proto/pbgo"
	mock "github.com/sktelecom/tks-proto/pbgo/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var contractID string

func getAccessor() (*contract.Accessor, error) {
	dsn := "host=localhost user=postgres password=password dbname=tks port=5432 sslmode=disable TimeZone=Asia/Seoul"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return contract.New(db), nil
}
func TestCreateContract(t *testing.T) {
	var err error
	s := server{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	mockInfoClient := mock.NewMockInfoServiceClient(ctrl)
	infoClient = gc.New(nil, mockInfoClient)
	contractAccessor, err = getAccessor()
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	defer cancel()

	contractorName := getRandomString("handler")
	t.Logf("New ContractorName %s", contractorName)

	req := pb.CreateContractRequest{
		ContractorName: contractorName,
		CspName:        "aws",
		CspAuth:        "{'token':'abcdefghijklmnop'}",
		Quota: &pb.ContractQuota{
			Cpu:    20,
			Memory: 40,
			Block:  12800000,
			Fs:     12800000,
		},
		AvailableServices: []string{"lma"},
	}
	// ctx, in.GetContractId(), in.GetCspName(), in.GetCspAuth())
	mockInfoClient.EXPECT().CreateCSPInfo(
		gomock.Any(),
		gomock.Any(),
	).Return(&pb.IDResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Id:    "a254a66e-7225-4527-bf4c-9b5494c99b37",
	}, nil)
	res, err := s.CreateContract(ctx, &req)
	if err != nil {
		t.Error("error occurred " + err.Error())
	}

	if res.Code != pb.Code_OK_UNSPECIFIED {
		t.Error("Not expected response code:", res.Code)
	}
	if res.CspId == "" {
		t.Error("CspId is empty.")
	}

	contractID = res.ContractId
}

func TestUpdateQuota(t *testing.T) {
	var err error
	s := server{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	contractAccessor, err = getAccessor()
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	req := pb.UpdateQuotaRequest{
		ContractId: contractID,
		Quota: &pb.ContractQuota{
			Cpu: 40,
		},
	}
	res, err := s.UpdateQuota(ctx, &req)
	if err != nil {
		t.Error("error occurred " + err.Error())
	}

	expected := &pb.UpdateQuotaResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		PrevQuota: &pb.ContractQuota{
			Cpu:    20,
			Memory: 40,
			Block:  12800000,
			Fs:     12800000,
		},
		CurrentQuota: &pb.ContractQuota{
			Cpu:    40,
			Memory: 40,
			Block:  12800000,
			Fs:     12800000,
		},
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("%v != %v", expected, res)
	}
}

func TestUpdateServices(t *testing.T) {
	var err error
	s := server{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	contractAccessor, err = getAccessor()
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	req := pb.UpdateServicesRequest{
		ContractId:        contractID,
		AvailableServices: []string{"lma", "servicemesh"},
	}
	res, err := s.UpdateServices(ctx, &req)
	if err != nil {
		t.Error("error occurred " + err.Error())
	}

	expected := &pb.UpdateServicesResponse{
		Code:            pb.Code_OK_UNSPECIFIED,
		Error:           nil,
		PrevServices:    []string{"lma"},
		CurrentServices: []string{"lma", "servicemesh"},
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("%v != %v", expected, res)
	}
}

func TestGetContract_InvalidArgument(t *testing.T) {
	s := server{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := pb.GetContractRequest{
		ContractId: "cbp-100002-xdkzkl",
	}
	res, _ := s.GetContract(ctx, &req)

	expected := &pb.GetContractResponse{
		Code: pb.Code_INVALID_ARGUMENT,
		Error: &pb.Error{
			Msg: "invalid contract ID cbp-100002-xdkzkl",
		},
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("%v != %v", expected, res)
	}
}

func TestGetQuota(t *testing.T) {
	s := server{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := pb.GetQuotaRequest{
		ContractId: contractID,
	}
	res, err := s.GetQuota(ctx, &req)
	if err != nil {
		t.Error("error occurred " + err.Error())
	}

	expected := &pb.GetQuotaResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Quota: &pb.ContractQuota{
			Cpu:      40,
			Memory:   40,
			Block:    12800000,
			BlockSsd: 0,
			Fs:       12800000,
			FsSsd:    0,
		},
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("%v != %v", expected, res)
	}
}

func TestGetAvailableServices(t *testing.T) {
	s := server{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := pb.GetAvailableServicesRequest{
		ContractId: contractID,
	}
	res, err := s.GetAvailableServices(ctx, &req)
	if err != nil {
		t.Error("error occurred " + err.Error())
	}

	expected := &pb.GetAvailableServicesResponse{
		Code:                pb.Code_OK_UNSPECIFIED,
		Error:               nil,
		AvaiableServiceApps: []string{"lma", "servicemesh"},
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("%v != %v", expected, res)
	}
}

func getRandomString(prefix string) string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return fmt.Sprintf("%s-%d", prefix, r.Int31n(1000000000))
}
