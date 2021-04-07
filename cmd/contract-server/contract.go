package main

import (
	"context"

	"github.com/openinfradev/tks-contract/pkg/log"
	pb "github.com/openinfradev/tks-proto/pbgo"
)

// CreateContract implements pbgo.ContractService.CreateContract gRPC
func (s *server) CreateContract(ctx context.Context, in *pb.CreateContractRequest) (*pb.CreateContractResponse, error) {
	log.Println("Not implemented: CreateContract")
	res := pb.CreateContractResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	return &res, nil
}

// UpdateQuota implements pbgo.ContractService.UpdateQuota gRPC
func (s *server) UpdateQuota(ctx context.Context, in *pb.UpdateQuotaRequest) (*pb.UpdateQuotaResponse, error) {
	log.Println("Not implemented: UpdateQuota")
	res := pb.UpdateQuotaResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	return &res, nil
}

// UpdateServices implements pbgo.ContractService.UpdateServices gRPC
func (s *server) UpdateServices(ctx context.Context, in *pb.UpdateServicesRequest) (*pb.UpdateServicesResponse, error) {
	log.Println("Not implemented: UpdateServices")
	res := pb.UpdateServicesResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	return &res, nil
}

// GetContract implements pbgo.ContractService.GetContract gRPC
func (s *server) GetContract(ctx context.Context, in *pb.GetContractRequest) (*pb.GetContractResponse, error) {
	log.Println("Not implemented: GetContract")
	res := pb.GetContractResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	return &res, nil
}

// GetQuota implements pbgo.ContractService.GetContract gRPC
func (s *server) GetQuota(ctx context.Context, in *pb.GetQuotaRequest) (*pb.GetQuotaResponse, error) {
	log.Println("Not implemented: GetQuota")
	res := pb.GetQuotaResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	return &res, nil
}

// GetServices implements pbgo.ContractService.GetServices gRPC
func (s *server) GetServices(ctx context.Context, in *pb.GetServicesRequest) (*pb.GetServicesResponse, error) {
	log.Println("Not implemented: GetServices")
	res := pb.GetServicesResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	return &res, nil
}
