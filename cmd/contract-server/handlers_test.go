package main

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/openinfradev/tks-proto/pbgo"
)

func TestCreateContract(t *testing.T) {
	s := server{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := pb.CreateContractRequest{
		ContractorName: "Tester",
		ContractId:     "cbp-100001-xdkzkl",
		Quota: &pb.ContractQuota{
			Cpu:    20,
			Memory: 40,
			Block:  12800000,
			Fs:     12800000,
		},
	}
	res, err := s.CreateContract(ctx, &req)
	if err != nil {
		t.Error("error occurred " + err.Error())
	}

	if res.Code != pb.Code_OK {
		t.Error("Not expected response code:", res.Code)
	}
	if res.McOpsId == "" {
		t.Error("McOpsId is empty.")
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
