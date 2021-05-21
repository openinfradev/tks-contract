package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	gc "github.com/sktelecom/tks-contract/pkg/grpc-client"
	pb "github.com/sktelecom/tks-proto/pbgo"
	mock "github.com/sktelecom/tks-proto/pbgo/mock"
)

func TestCreateContract(t *testing.T) {
	s := server{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	mockInfoClient := mock.NewMockInfoServiceClient(ctrl)
	infoClient = gc.New(nil, mockInfoClient)
	defer cancel()
	req := pb.CreateContractRequest{
		ContractorName: "Tester",
		ContractId:     "cbp-100001-xdkzkl",
		CspName:        "aws",
		CspAuth:        "{'token':'abcdefghijklmnop'}",
		Quota: &pb.ContractQuota{
			Cpu:    20,
			Memory: 40,
			Block:  12800000,
			Fs:     12800000,
		},
	}
	// ctx, in.GetContractId(), in.GetCspName(), in.GetCspAuth())
	mockInfoClient.EXPECT().CreateCSPInfo(
		gomock.Any(),
		&pb.CreateCSPInfoRequest{
			ContractId: req.ContractId,
			CspName:    req.CspName,
			Auth:       req.CspAuth,
		},
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
}

func TestUpdateQuota(t *testing.T) {
	s := server{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := pb.UpdateQuotaRequest{
		ContractId: "cbp-100001-xdkzkl",
		Quota: &pb.ContractQuota{
			Cpu:    20,
			Memory: 40,
			Block:  12800000,
			Fs:     12800000,
		},
	}
	res, err := s.UpdateQuota(ctx, &req)
	if err != nil {
		t.Error("error occurred " + err.Error())
	}

	expected := &pb.UpdateQuotaResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("%v != %v", expected, res)
	}
}

func TestUpdateServices(t *testing.T) {
	s := server{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := pb.UpdateServicesRequest{
		ContractId:        "cbp-100001-xdkzkl",
		AvailableServices: []string{"lma"},
	}
	res, err := s.UpdateServices(ctx, &req)
	if err != nil {
		t.Error("error occurred " + err.Error())
	}

	expected := &pb.UpdateServicesResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("%v != %v", expected, res)
	}
}

func TestGetContract(t *testing.T) {
	s := server{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := pb.GetContractRequest{
		ContractId: "cbp-100002-xdkzkl",
	}
	res, _ := s.GetContract(ctx, &req)

	expected := &pb.GetContractResponse{
		Code: pb.Code_NOT_FOUND,
		Error: &pb.Error{
			Msg: "Not found contract for cbp-100002-xdkzkl",
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
		ContractId: "cbp-100001-xdkzkl",
	}
	res, err := s.GetQuota(ctx, &req)
	if err != nil {
		t.Error("error occurred " + err.Error())
	}

	expected := &pb.GetQuotaResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("%v != %v", expected, res)
	}
}
