package main

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
	"errors"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/stretchr/testify/require"
	"github.com/golang/mock/gomock"
	
	"github.com/openinfradev/tks-common/pkg/log"
	"github.com/openinfradev/tks-common/pkg/helper"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	mock "github.com/openinfradev/tks-proto/tks_pb/mock"
	
	"github.com/openinfradev/tks-contract/pkg/contract"
	model "github.com/openinfradev/tks-contract/pkg/contract/model"

	mockargo "github.com/openinfradev/tks-common/pkg/argowf/mock"
)



var (
	err error
	testDBHost string
	testDBPort string
)

var (
	contractId string
	requestForSenariTest = randomRequest()
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

func TestCreateContract(t *testing.T){
	testCases := []struct {
		name			string
		in				pb.CreateContractRequest
		buildStubs		func(mockInfoClient *mock.MockCspInfoServiceClient, mockArgoClient *mockargo.MockClient)
		checkResponse	func(req *pb.CreateContractRequest, res *pb.CreateContractResponse, err error )
	}{
		{
			name: "OK",
			in: requestForSenariTest,
			buildStubs: func(mockInfoClient *mock.MockCspInfoServiceClient, mockArgoClient *mockargo.MockClient) {
				mockInfoClient.EXPECT().CreateCSPInfo( gomock.Any(), gomock.Any(), ).Return(&pb.IDResponse{
					Code:  pb.Code_OK_UNSPECIFIED,
					Error: nil,
					Id:    uuid.New().String(),
				}, nil)

				mockArgoClient.EXPECT().
					SumbitWorkflowFromWftpl(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any() ).
					Times(1).
					Return(randomString("workflowName"), nil)
			},
			checkResponse: func(req *pb.CreateContractRequest, res *pb.CreateContractResponse, err error){
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
				require.True(t, res.ContractId != "" )

				// store for senario test
				requestForSenariTest = *req
				contractId = res.ContractId
			},
		},
		{
			name: "NOT_FOUND",
			in: requestForSenariTest,
			buildStubs: func(mockInfoClient *mock.MockCspInfoServiceClient, mockArgoClient *mockargo.MockClient) {
			},
			checkResponse: func(req *pb.CreateContractRequest, res *pb.CreateContractResponse, err error){
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_NOT_FOUND)
			},
		},
		{
			name: "CSP_INVALID_ARGUMENT",
			in: randomRequest(),
			buildStubs: func(mockInfoClient *mock.MockCspInfoServiceClient, mockArgoClient *mockargo.MockClient) {
				id := uuid.New()
				mockInfoClient.EXPECT().CreateCSPInfo( gomock.Any(), gomock.Any(), ).
					Return(&pb.IDResponse{
						Code:  pb.Code_INVALID_ARGUMENT,
						Error: &pb.Error{
	       					Msg: fmt.Sprintf("invalid contract ID %s", id),
	       				},
					}, nil)
			},
			checkResponse: func(req *pb.CreateContractRequest, res *pb.CreateContractResponse, err error){
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "ARGO_WORKFLOW_ERROR",
			in: randomRequest(),
			buildStubs: func(mockInfoClient *mock.MockCspInfoServiceClient, mockArgoClient *mockargo.MockClient) {
				id := uuid.New()
				mockInfoClient.EXPECT().CreateCSPInfo( gomock.Any(), gomock.Any(), ).Return(&pb.IDResponse{
					Code:  pb.Code_OK_UNSPECIFIED,
					Error: nil,
					Id:    id.String(),
				}, nil)

				mockArgoClient.EXPECT().
					SumbitWorkflowFromWftpl(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any() ).
					Times(1).
					Return(randomString("workflowName"), errors.New("argo error"))
			},
			checkResponse: func(req *pb.CreateContractRequest, res *pb.CreateContractResponse, err error){
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_INTERNAL)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			contractAccessor, err = getAccessor()

			// mocking and injection
			mockArgoClient := mockargo.NewMockClient(ctrl)
			argowfClient = mockArgoClient 
			mockInfoClient := mock.NewMockCspInfoServiceClient(ctrl)
			cspInfoClient = mockInfoClient
			
			tc.buildStubs(mockInfoClient, mockArgoClient)

			s := server{}
			res, err := s.CreateContract(ctx, &tc.in)
			tc.checkResponse(&tc.in, res, err)
		})
	}

}

func TestUpdateQuota(t *testing.T){
	testCases := []struct {
		name			string
		in				pb.UpdateQuotaRequest
		checkResponse	func(req *pb.UpdateQuotaRequest, res *pb.UpdateQuotaResponse, err error )
	}{
		{
			name: "OK",
			in: pb.UpdateQuotaRequest{
				ContractId: contractId,
				Quota: &pb.ContractQuota{
					Cpu: 40,
				},
			},
			checkResponse: func(req *pb.UpdateQuotaRequest, res *pb.UpdateQuotaResponse, err error){
				expected := &pb.UpdateQuotaResponse{
					Code:  pb.Code_OK_UNSPECIFIED,
					Error: nil,
					PrevQuota: &pb.ContractQuota{
						Cpu:    20,
						Memory: 40,
						Block:  12800000,
						BlockSsd: 12800000,
						Fs:     12800000,
						FsSsd:  12800000,
					},
					CurrentQuota: &pb.ContractQuota{
						Cpu:    40,
						Memory: 40,
						Block:  12800000,
						Fs:     12800000,
						BlockSsd: 12800000,
						FsSsd:  12800000,
					},
				}

				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
				require.Equal(t, reflect.DeepEqual(expected, res), true)

				requestForSenariTest.Quota.Cpu = 40
			},
		},
		{
			name: "RECORD_NOT_FOUND",
			in: pb.UpdateQuotaRequest{
				ContractId: uuid.New().String(),
				Quota: &pb.ContractQuota{
					Cpu: 40,
				},
			},
			checkResponse: func(req *pb.UpdateQuotaRequest, res *pb.UpdateQuotaResponse, err error){
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INTERNAL)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			contractAccessor, err = getAccessor()

			s := server{}
			res, err := s.UpdateQuota(ctx, &tc.in)

			tc.checkResponse(&tc.in, res, err)
		})
	}

}

func TestUpdateServices(t *testing.T){
	testCases := []struct {
		name			string
		in				pb.UpdateServicesRequest
		checkResponse	func(req *pb.UpdateServicesRequest, res *pb.UpdateServicesResponse, err error )
	}{
		{
			name: "OK",
			in: pb.UpdateServicesRequest{
				ContractId:        contractId,
				AvailableServices: []string{"lma", "servicemesh"},
			},
			checkResponse: func(req *pb.UpdateServicesRequest, res *pb.UpdateServicesResponse, err error){
				expected := &pb.UpdateServicesResponse{
					Code:            pb.Code_OK_UNSPECIFIED,
					Error:           nil,
					PrevServices:    []string{"lma"},
					CurrentServices: []string{"lma", "servicemesh"},
				}

				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
				require.Equal(t, reflect.DeepEqual(expected, res), true)
			},
		},
		{
			name: "OK_DELETE_SERVICE",
			in: pb.UpdateServicesRequest{
				ContractId:        contractId,
				AvailableServices: []string{""},
			},
			checkResponse: func(req *pb.UpdateServicesRequest, res *pb.UpdateServicesResponse, err error){
				expected := &pb.UpdateServicesResponse{
					Code:            pb.Code_OK_UNSPECIFIED,
					Error:           nil,
					PrevServices:    []string{"lma", "servicemesh"},
					CurrentServices: []string{""},
				}

				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
				require.Equal(t, reflect.DeepEqual(expected, res), true)
				requestForSenariTest.AvailableServices = []string{""};
			},
		},
		{
			name: "RECORD_NOT_FOUND",
			in: pb.UpdateServicesRequest{
				ContractId: uuid.New().String(),
				AvailableServices: []string{"lma", "servicemesh"},
			},
			checkResponse: func(req *pb.UpdateServicesRequest, res *pb.UpdateServicesResponse, err error){
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INTERNAL)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			contractAccessor, err = getAccessor()

			s := server{}
			res, err := s.UpdateServices(ctx, &tc.in)

			tc.checkResponse(&tc.in, res, err)
		})
	}

}

func TestGetContract(t *testing.T){
	testCases := []struct {
		name			string
		in				pb.GetContractRequest
		checkResponse	func(req *pb.GetContractRequest, res *pb.GetContractResponse, err error )
	}{
		{
			name: "OK",
			in: pb.GetContractRequest{
				ContractId: contractId,
			},
			checkResponse: func(req *pb.GetContractRequest, res *pb.GetContractResponse, err error){
				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)

				contract := res.GetContract()
				require.NotNil(t, contract)
				require.Equal(t, contractId, contract.GetContractId())
				require.Equal(t, requestForSenariTest.ContractorName, contract.GetContractorName())
			},
		},
		{
			name: "INVALID_CONTRACT_ID",
			in: pb.GetContractRequest{
				ContractId: "invalid_contract_id",
			},
			checkResponse: func(req *pb.GetContractRequest, res *pb.GetContractResponse, err error){
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOT_FOUND",
			in: pb.GetContractRequest{
				ContractId: uuid.New().String(),
			},
			checkResponse: func(req *pb.GetContractRequest, res *pb.GetContractResponse, err error){
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_NOT_FOUND)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			contractAccessor, err = getAccessor()

			s := server{}
			res, err := s.GetContract(ctx, &tc.in)

			tc.checkResponse(&tc.in, res, err)
		})
	}

}

func TestGetQuota(t *testing.T){
	testCases := []struct {
		name			string
		in				pb.GetQuotaRequest
		checkResponse	func(req *pb.GetQuotaRequest, res *pb.GetQuotaResponse, err error )
	}{
		{
			name: "OK",
			in: pb.GetQuotaRequest{
				ContractId: contractId,
			},
			checkResponse: func(req *pb.GetQuotaRequest, res *pb.GetQuotaResponse, err error){
				expected := &pb.GetQuotaResponse{
					Code:  pb.Code_OK_UNSPECIFIED,
					Error: nil,
					Quota: &pb.ContractQuota{
						Cpu:      requestForSenariTest.Quota.Cpu,
						Memory:   requestForSenariTest.Quota.Memory,
						Block:    requestForSenariTest.Quota.Block,
						BlockSsd: requestForSenariTest.Quota.BlockSsd,
						Fs:       requestForSenariTest.Quota.Fs,
						FsSsd:    requestForSenariTest.Quota.FsSsd,
					},
				}

				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
				require.Equal(t, reflect.DeepEqual(expected, res), true)
			},
		},
		{
			name: "INVALID_CONTRACT_ID",
			in: pb.GetQuotaRequest{
				ContractId: "invalid_contract_id",
			},
			checkResponse: func(req *pb.GetQuotaRequest, res *pb.GetQuotaResponse, err error){
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOT_FOUND",
			in: pb.GetQuotaRequest{
				ContractId: uuid.New().String(),
			},
			checkResponse: func(req *pb.GetQuotaRequest, res *pb.GetQuotaResponse, err error){
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_NOT_FOUND)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			contractAccessor, err = getAccessor()

			s := server{}
			res, err := s.GetQuota(ctx, &tc.in)

			tc.checkResponse(&tc.in, res, err)
		})
	}

}

func TestGetAvailableServices(t *testing.T){
	testCases := []struct {
		name			string
		in				pb.GetAvailableServicesRequest
		checkResponse	func(req *pb.GetAvailableServicesRequest, res *pb.GetAvailableServicesResponse, err error )
	}{
		{
			name: "OK",
			in: pb.GetAvailableServicesRequest{
				ContractId: contractId,
			},
			checkResponse: func(req *pb.GetAvailableServicesRequest, res *pb.GetAvailableServicesResponse, err error){
				expected := &pb.GetAvailableServicesResponse{
					Code:                pb.Code_OK_UNSPECIFIED,
					Error:               nil,
					AvaiableServiceApps: requestForSenariTest.AvailableServices,
				}

				require.NoError(t, err)
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
				require.Equal(t, reflect.DeepEqual(expected, res), true)


			},
		},
		{
			name: "INVALID_CONTRACT_ID",
			in: pb.GetAvailableServicesRequest{
				ContractId: "invalid_contract_id",
			},
			checkResponse: func(req *pb.GetAvailableServicesRequest, res *pb.GetAvailableServicesResponse, err error){
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "NOT_FOUND",
			in: pb.GetAvailableServicesRequest{
				ContractId: uuid.New().String(),
			},
			checkResponse: func(req *pb.GetAvailableServicesRequest, res *pb.GetAvailableServicesResponse, err error){
				require.Error(t, err)
				require.Equal(t, res.Code, pb.Code_NOT_FOUND)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			contractAccessor, err = getAccessor()

			s := server{}
			res, err := s.GetAvailableServices(ctx, &tc.in)

			tc.checkResponse(&tc.in, res, err)
		})
	}

}


// Helpers

func randomString(prefix string) string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return fmt.Sprintf("%s-%d", prefix, r.Int31n(1000000000))
}

func randomContract() (model.Contract) {
	return model.Contract {
		ID: uuid.New(),
		ContractorName: randomString("NAME"),
		AvailableServices: []string{"lma"},
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

func randomResourceQuota() (model.ResourceQuota) {
	return model.ResourceQuota {
		ID: uuid.New(),
		Cpu: 20,
		Memory: 40,
		Block: 12800000,
		BlockSsd: 12800000,
		Fs: 12800000,
		FsSsd: 12800000,
		//ContractID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func reflectToRequest( contract model.Contract, resourceQuota model.ResourceQuota ) (pb.CreateContractRequest) {
	return pb.CreateContractRequest {
		ContractorName: contract.ContractorName,
		CspName: "aws",
		CspAuth: "{'token':'csp_auth_token'}",
		AvailableServices: contract.AvailableServices,
		Quota: &pb.ContractQuota{
			Cpu: resourceQuota.Cpu,
			Memory: resourceQuota.Memory,
			Block: resourceQuota.Block,
			BlockSsd: resourceQuota.BlockSsd,
			Fs: resourceQuota.Fs,
			FsSsd: resourceQuota.FsSsd,
		},
	}
}

func randomRequest() (pb.CreateContractRequest){
	testContract := randomContract()
	testResourceQuota := randomResourceQuota()
	return reflectToRequest(testContract, testResourceQuota)
}
