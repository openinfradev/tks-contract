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
		t.Error("error occured " + err.Error())
	}

	expected := &pb.CreateContractResponse{
		Code:    pb.Code_UNIMPLEMENTED,
		Error:   nil,
		McOpsId: "",
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("%v != %v", expected, res)
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
		t.Error("error occured " + err.Error())
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
		t.Error("error occured " + err.Error())
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
		ContractId: "cbp-100001-xdkzkl",
	}
	res, err := s.GetContract(ctx, &req)
	if err != nil {
		t.Error("error occured " + err.Error())
	}

	expected := &pb.GetContractResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
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
		t.Error("error occured " + err.Error())
	}

	expected := &pb.GetQuotaResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("%v != %v", expected, res)
	}
}