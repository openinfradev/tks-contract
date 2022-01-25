package main

import (
	"context"
	"testing"
	"math/rand"
	"time"
	"fmt"
	"errors"

	"github.com/google/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	accessor "github.com/openinfradev/tks-contract/pkg/contract"
	contract "github.com/openinfradev/tks-contract/pkg/contract/model"
	gc "github.com/openinfradev/tks-contract/pkg/grpc-client"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	mock "github.com/openinfradev/tks-proto/tks_pb/mock"
	mockargo "github.com/openinfradev/tks-cluster-lcm/pkg/argowf/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"database/sql"

)

type dbConn struct {
	db      *gorm.DB
	mock    sqlmock.Sqlmock
}

func getAccessor() (*accessor.Accessor, sqlmock.Sqlmock, error) {
	s := &dbConn{}
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	s.db, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return accessor.New(s.db), s.mock, nil
}

func TestCreateContract(t *testing.T) {
	contract := randomContract();
	resourceQuota := randomResourceQuota(contract.ID);
	req := reflectToRequest(contract, resourceQuota)

	testCases := []struct {
		name			string
		in				pb.CreateContractRequest
		buildStubs		func( mockSql sqlmock.Sqlmock, mockInfoClient *mock.MockCspInfoServiceClient, mockArgoClient *mockargo.MockClient )
		checkResponse	func(res *pb.CreateContractResponse)
	}{
		{
			//("contractor_name","available_services","updated_at","created_at","id"
			name: "OK",
			in: req,
			buildStubs: func( mockSql sqlmock.Sqlmock, mockInfoClient *mock.MockCspInfoServiceClient, mockArgoClient *mockargo.MockClient ) {
				mockSql.ExpectBegin()
				mockSql.ExpectQuery(`INSERT INTO "contracts" (.+) RETURNING`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
				mockSql.ExpectQuery(`INSERT INTO "resource_quota" (.+) RETURNING`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
				mockSql.ExpectCommit()

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
			checkResponse: func(res *pb.CreateContractResponse){
				require.Equal(t, res.Code, pb.Code_OK_UNSPECIFIED)
			},
		},
		{
			// [TODO] NOT_FOUND 가 아니라, DB_ERROR 로 변경해야 할 듯함.
			name: "DB_ERROR_INTERNAL",
			in: req,
			buildStubs: func( mockSql sqlmock.Sqlmock, mockInfoClient *mock.MockCspInfoServiceClient, mockArgoClient *mockargo.MockClient ) {
				mockSql.ExpectBegin()
				mockSql.ExpectQuery(`INSERT INTO "contracts" (.+) RETURNING`).
					WillReturnError(fmt.Errorf("some error"))
				mockSql.ExpectQuery(`INSERT INTO "resource_quota" (.+) RETURNING`).
					WillReturnError(fmt.Errorf("some error"))
				mockSql.ExpectRollback()
			},
			checkResponse: func(res *pb.CreateContractResponse){
				require.Equal(t, res.Code, pb.Code_NOT_FOUND)
			},
		},
		{
			name: "CSP_INVALID_ARGUMENT",
			in: req,
			buildStubs: func( mockSql sqlmock.Sqlmock, mockInfoClient *mock.MockCspInfoServiceClient, mockArgoClient *mockargo.MockClient ) {
				id := uuid.New()
				mockSql.ExpectBegin()
				mockSql.ExpectQuery(`INSERT INTO "contracts" (.+) RETURNING`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
				mockSql.ExpectQuery(`INSERT INTO "resource_quota" (.+) RETURNING`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
				mockSql.ExpectCommit()

				mockInfoClient.EXPECT().CreateCSPInfo( gomock.Any(), gomock.Any(), ).
					Return(&pb.IDResponse{
						Code:  pb.Code_INVALID_ARGUMENT,
						Error: &pb.Error{
	       					Msg: fmt.Sprintf("invalid contract ID %s", id),
	       				},
					}, nil)
			},
			checkResponse: func(res *pb.CreateContractResponse){
				require.Equal(t, res.Code, pb.Code_INVALID_ARGUMENT)
			},
		},
		{
			name: "ARGO_WORKFLOW_ERROR",
			in: req,
			buildStubs: func( mockSql sqlmock.Sqlmock, mockInfoClient *mock.MockCspInfoServiceClient, mockArgoClient *mockargo.MockClient ) {
				id := uuid.New()
				mockSql.ExpectBegin()
				mockSql.ExpectQuery(`INSERT INTO "contracts" (.+) RETURNING`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
				mockSql.ExpectQuery(`INSERT INTO "resource_quota" (.+) RETURNING`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
				mockSql.ExpectCommit()

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
			checkResponse: func(res *pb.CreateContractResponse){
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

			mockAccessor, mockSql, _ := getAccessor()
			mockInfoClient := mock.NewMockCspInfoServiceClient(ctrl)
			mockArgoClient := mockargo.NewMockClient(ctrl)

			// injection
			contractAccessor = mockAccessor
			cspInfoClient = gc.NewCspInfoServiceClient(nil, mockInfoClient)
			argowfClient = mockArgoClient 
			
			tc.buildStubs(mockSql, mockInfoClient, mockArgoClient)

			s := server{}
			res, err := s.CreateContract(ctx, &tc.in)
			require.NoError(t, err)

			tc.checkResponse(res)
		})
	}
}



type ResourceQuota struct {
	ID         uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Cpu        int64
	Memory     int64
	Block      int64
	BlockSsd   int64
	Fs         int64
	FsSsd      int64
	ContractID uuid.UUID `gorm:"type:uuid"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func randomContract() (contract.Contract) {
	return contract.Contract {
		ID: uuid.New(),
		ContractorName: randomString("NAME"),
		AvailableServices: []string{"lma"},
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

func randomResourceQuota( contractId uuid.UUID) (contract.ResourceQuota) {
	return contract.ResourceQuota {
		ID: uuid.New(),
		Cpu: 20,
		Memory: 40,
		Block: 12800000,
		BlockSsd: 12800000,
		Fs: 12800000,
		FsSsd: 12800000,
		ContractID: contractId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func reflectToRequest( contract contract.Contract, resourceQuota contract.ResourceQuota ) (pb.CreateContractRequest) {
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

func randomString(prefix string) string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return fmt.Sprintf("%s-%d", prefix, r.Int31n(1000000000))
}
