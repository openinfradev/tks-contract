package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/openinfradev/tks-common/pkg/log"
	"github.com/openinfradev/tks-common/pkg/argowf"
	pb "github.com/openinfradev/tks-proto/tks_pb"

	"github.com/openinfradev/tks-contract/pkg/contract"
)

var (
	argowfClient argowf.Client
	contractAccessor *contract.Accessor
	cspInfoClient    pb.CspInfoServiceClient
)

func InitHandlers( argoAddress string, argoPort int ) {
	_client, err := argowf.New( argoAddress, argoPort );
	if err != nil {
		log.Fatal( "failed to create argowf client : ", err )
	}
	argowfClient = _client;
}

// CreateContract implements pbgo.ContractService.CreateContract gRPC
func (s *server) CreateContract(ctx context.Context, in *pb.CreateContractRequest) (*pb.CreateContractResponse, error) {
	log.Info("Request 'CreateContract' for contract name", in.GetContractorName())

	contractId, err := contractAccessor.Create(in.GetContractorName(), in.GetAvailableServices(), in.GetQuota())
	if err != nil {
		return &pb.CreateContractResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, nil
	}
	log.Info("newly created Contract Id:", contractId)

	res, err := cspInfoClient.CreateCSPInfo(ctx, &pb.CreateCSPInfoRequest{ 
		ContractId: contractId.String(), 
		CspName: in.GetCspName(), 
		Auth: in.GetCspAuth(),
	})
	log.Info("newly created CSP Id:", res.GetId())
	if err != nil {
		return &pb.CreateContractResponse{
			Code: res.GetCode(),
			Error: &pb.Error{
				Msg: err.Error(),
			}, 
		}, nil
	}

	if res.GetCode() != pb.Code_OK_UNSPECIFIED {
		return &pb.CreateContractResponse{
			Code: res.GetCode(),
			Error: &pb.Error{
				Msg: res.GetError().String(),
			}, 
		}, nil
	}

	workflowTemplate := "tks-create-contract-repo"
	nameSpace := "argo"
	parameters := []string{ 
		"contract_id=" + contractId.String(),  
	};

	workflowName, err := argowfClient.SumbitWorkflowFromWftpl( ctx, workflowTemplate, nameSpace, parameters );
	if err != nil {
		log.Error( "failed to submit argo workflow template. err : ", err )
		return &pb.CreateContractResponse {
			Code: pb.Code_INTERNAL,
			Error: &pb.Error{
				Msg: fmt.Sprintf("Failed to call argo workflow : %s", err ),
			},
		}, nil
	}
	log.Info("submited workflow :", workflowName )

	//argowfClient.WaitWorkflows(ctx, nameSpace, []string{workflowName}, false, false)

	//log.Info("completed workflow :", workflowName )

	return &pb.CreateContractResponse{
		Code:       pb.Code_OK_UNSPECIFIED,
		Error:      nil,
		CspId:      res.GetId(),
		ContractId: contractId.String(),
	}, nil
}

// UpdateQuota implements pbgo.ContractService.UpdateQuota gRPC
func (s *server) UpdateQuota(ctx context.Context, in *pb.UpdateQuotaRequest) (*pb.UpdateQuotaResponse, error) {
	log.Info("Request 'UpdateQuota' for contract id ", in.GetContractId())
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
	prev, curr, err := contractAccessor.UpdateResourceQuota(contractID, in.GetQuota())

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
	log.Info("Request 'UpdateServices' for contract id ", in.GetContractId())
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
	log.Info("Request 'GetContract' for contract id ", in.GetContractId())
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
	log.Info("Request 'GetQuota' for contract id ", in.GetContractId())
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
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: err.Error(),
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

// GetAvailableServices implements pbgo.ContractService.GetAvailableServices gRPC
func (s *server) GetAvailableServices(ctx context.Context, in *pb.GetAvailableServicesRequest) (*pb.GetAvailableServicesResponse, error) {
	log.Info("Request 'GetAvailableServices' for contract id ", in.GetContractId())
	contractID, err := uuid.Parse(in.GetContractId())
	if err != nil {
		return &pb.GetAvailableServicesResponse{
			Code: pb.Code_INVALID_ARGUMENT,
			Error: &pb.Error{
				Msg: fmt.Sprintf("invalid contract ID %s", in.GetContractId()),
			},
		}, err
	}

	contract, err := contractAccessor.GetContract(contractID)
	if err != nil {
		return &pb.GetAvailableServicesResponse{
			Code: pb.Code_NOT_FOUND,
			Error: &pb.Error{
				Msg: err.Error(),
			},
		}, err
	}

	res := pb.GetAvailableServicesResponse{
		Code:                pb.Code_OK_UNSPECIFIED,
		Error:               nil,
		AvaiableServiceApps: contract.AvailableServices,
	}
	return &res, nil
}
