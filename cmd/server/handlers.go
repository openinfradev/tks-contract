package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sktelecom/tks-contract/pkg/contract"
	gc "github.com/sktelecom/tks-contract/pkg/grpc-client"
	"github.com/sktelecom/tks-contract/pkg/log"
	pb "github.com/sktelecom/tks-proto/pbgo"
)

var (
	contractAccessor *contract.Accessor
	infoClient       *gc.InfoClient
)

// CreateContract implements pbgo.ContractService.CreateContract gRPC
func (s *server) CreateContract(ctx context.Context, in *pb.CreateContractRequest) (*pb.CreateContractResponse, error) {
	log.Info("Request 'CreateContract' for contract name", in.GetContractorName())
	inputQuota := in.GetQuota()
	contractID, err := contractAccessor.Create(in.GetContractorName(),
		in.GetAvailableServices(),
		contract.ResourceQuotaParam{
			Cpu:      inputQuota.Cpu,
			Memory:   inputQuota.Memory,
			Block:    inputQuota.Block,
			BlockSsd: inputQuota.BlockSsd,
			Fs:       inputQuota.Fs,
			FsSsd:    inputQuota.FsSsd,
		})
	if err != nil {
		res := pb.CreateContractResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}
		return &res, err
	}
	// Create New CSP Info
	res, err := infoClient.CreateCSPInfo(ctx, contractID.String(), in.GetCspName(), in.GetCspAuth())
	log.Info("newly created CSP ID:", res.GetId())
	if err != nil || res.GetCode() != pb.Code_OK_UNSPECIFIED {
		res := pb.CreateContractResponse{
			Code: res.GetCode(),
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}
		return &res, err
	}
	return &pb.CreateContractResponse{
		Code:       pb.Code_OK_UNSPECIFIED,
		Error:      nil,
		CspId:      res.GetId(),
		ContractId: contractID.String(),
	}, nil
}

// UpdateQuota implements pbgo.ContractService.UpdateQuota gRPC
func (s *server) UpdateQuota(ctx context.Context, in *pb.UpdateQuotaRequest) (*pb.UpdateQuotaResponse, error) {
	log.Info("Request 'UpdateQuota' for contract id", in.GetContractId())
	contractID, err := uuid.Parse(in.GetContractId())
	if err != nil {
		res := pb.UpdateQuotaResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid contract ID %s", in.GetContractId()),
			},
		}
		return &res, err
	}
	inputQuota := in.GetQuota()
	prev, curr, err := contractAccessor.UpdateResourceQuota(contractID, contract.ResourceQuotaParam{
		Cpu:      inputQuota.Cpu,
		Memory:   inputQuota.Memory,
		Block:    inputQuota.Block,
		BlockSsd: inputQuota.BlockSsd,
		Fs:       inputQuota.Fs,
		FsSsd:    inputQuota.FsSsd,
	})

	if err != nil {
		res := pb.UpdateQuotaResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}
		return &res, err
	}
	return &pb.UpdateQuotaResponse{
		Code:         pb.Code_OK_UNSPECIFIED,
		Error:        nil,
		PrevQuota:    prev,
		CurrentQuota: curr,
	}, nil
}

// UpdateServices implements pbgo.ContractService.UpdateServices gRPC
func (s *server) UpdateServices(ctx context.Context, in *pb.UpdateServicesRequest) (*pb.UpdateServicesResponse, error) {
	log.Info("Request 'UpdateServices' for contract id", in.GetContractId())
	contractID, err := uuid.Parse(in.GetContractId())
	if err != nil {
		res := pb.UpdateServicesResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid contract ID %s", in.GetContractId()),
			},
		}
		return &res, err
	}
	prev, curr, err := contractAccessor.UpdateAvailableServices(contractID, in.GetAvailableServices())
	if err != nil {
		res := pb.UpdateServicesResponse{
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}
		return &res, err
	}
	return &pb.UpdateServicesResponse{
		Code:            pb.Code_OK_UNSPECIFIED,
		Error:           nil,
		PrevServices:    prev,
		CurrentServices: curr,
	}, nil
}

// GetContract implements pbgo.ContractService.GetContract gRPC
func (s *server) GetContract(ctx context.Context, in *pb.GetContractRequest) (*pb.GetContractResponse, error) {
	log.Info("Request 'GetContract' for contractID", in.GetContractId())
	contractID, err := uuid.Parse(in.GetContractId())
	if err != nil {
		res := pb.GetContractResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid contract ID %s", in.GetContractId()),
			},
		}
		return &res, err
	}
	contract, err := contractAccessor.GetContract(contractID)
	if err != nil {
		res := pb.GetContractResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}
		return &res, err
	}
	res := pb.GetContractResponse{
		Code:     pb.Code_OK_UNSPECIFIED,
		Error:    nil,
		Contract: contract,
	}
	return &res, nil
}

// GetQuota implements pbgo.ContractService.GetContract gRPC
func (s *server) GetQuota(ctx context.Context, in *pb.GetQuotaRequest) (*pb.GetQuotaResponse, error) {
	log.Info("Request 'GetQuota' for contractID", in.GetContractId())
	contractID, err := uuid.Parse(in.GetContractId())
	if err != nil {
		return &pb.GetQuotaResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid contract ID %s", in.GetContractId()),
			},
		}, err
	}

	quota, err := contractAccessor.GetResourceQuota(contractID)
	if err != nil {
		return &pb.GetQuotaResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid contract ID %s", in.GetContractId()),
			},
		}, err
	}
	return &pb.GetQuotaResponse{
		Code:  pb.Code_OK_UNSPECIFIED,
		Error: nil,
		Quota: &pb.ContractQuota{
			Cpu:      quota.Cpu,
			Memory:   quota.Memory,
			Block:    quota.Block,
			BlockSsd: quota.BlockSsd,
			Fs:       quota.Fs,
			FsSsd:    quota.FsSsd,
		},
	}, nil
}

// GetServices implements pbgo.ContractService.GetServices gRPC
func (s *server) GetServices(ctx context.Context, in *pb.GetServicesRequest) (*pb.GetServicesResponse, error) {
	log.Warn("Not implemented: GetServices")
	res := pb.GetServicesResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	return &res, nil
}
