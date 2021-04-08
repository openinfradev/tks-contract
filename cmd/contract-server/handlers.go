package main

import (
	"context"

	"github.com/openinfradev/tks-contract/pkg/contract"
	"github.com/openinfradev/tks-contract/pkg/log"
	pb "github.com/openinfradev/tks-proto/pbgo"
)

var contractAccessor *contract.ContractAccessor

// CreateContract implements pbgo.ContractService.CreateContract gRPC
func (s *server) CreateContract(ctx context.Context, in *pb.CreateContractRequest) (*pb.CreateContractResponse, error) {
	log.Debug("Request 'CreateContract' for contractID", in.GetContractId())
	mID, err := contractAccessor.Post(in.GetContractorName(),
		contract.ContractId(in.GetContractId()),
		in.GetAvailableServices(),
		in.GetQuota())
	if err != nil {
		res := pb.CreateContractResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}
		return &res, err
	}
	res := pb.CreateContractResponse{
		Code:    pb.Code_OK,
		Error:   nil,
		McOpsId: mID.String(),
	}
	return &res, nil
}

// UpdateQuota implements pbgo.ContractService.UpdateQuota gRPC
func (s *server) UpdateQuota(ctx context.Context, in *pb.UpdateQuotaRequest) (*pb.UpdateQuotaResponse, error) {
	log.Debug("Not implemented: UpdateQuota")
	res := pb.UpdateQuotaResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	return &res, nil
}

// UpdateServices implements pbgo.ContractService.UpdateServices gRPC
func (s *server) UpdateServices(ctx context.Context, in *pb.UpdateServicesRequest) (*pb.UpdateServicesResponse, error) {
	log.Debug("Not implemented: UpdateServices")
	res := pb.UpdateServicesResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	return &res, nil
}

// GetContract implements pbgo.ContractService.GetContract gRPC
func (s *server) GetContract(ctx context.Context, in *pb.GetContractRequest) (*pb.GetContractResponse, error) {
	log.Debug("Request 'GetContract' for contractID", in.GetContractId())
	doc, err := contractAccessor.Get(contract.ContractId(in.GetContractId()))
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
		Code:              pb.Code_OK,
		Error:             nil,
		ContractorName:    doc.ContractorName,
		ContractId:        string(doc.ContractId),
		Quota:             doc.Quota,
		AvailableServices: doc.AvailableServices,
		McOpsId:           doc.McOpsId.String(),
		LastUpdatedTs:     doc.LastUpdatedTs.Timestamppb(),
	}
	return &res, nil
}

// GetQuota implements pbgo.ContractService.GetContract gRPC
func (s *server) GetQuota(ctx context.Context, in *pb.GetQuotaRequest) (*pb.GetQuotaResponse, error) {
	log.Warn("Not implemented: GetQuota")
	res := pb.GetQuotaResponse{
		Code:  pb.Code_UNIMPLEMENTED,
		Error: nil,
	}
	return &res, nil
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
