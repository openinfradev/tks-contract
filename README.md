# tks-contract

[![Go Report Card](https://goreportcard.com/badge/github.com/openinfradev/tks-contract?style=flat-square)](https://goreportcard.com/report/github.com/openinfradev/tks-contract)
[![Go Reference](https://pkg.go.dev/badge/github.com/openinfradev/tks-contract.svg)](https://pkg.go.dev/github.com/openinfradev/tks-contract)
[![Release](https://img.shields.io/github/release/sktelecom/tks-contract.svg?style=flat-square)](https://github.com/openinfradev/tks-contract/releases/latest)

TKS는 Taco Kubernetes Service의 약자로, SK Telecom이 만든 GitOps기반의 서비스 시스템을 의미합니다. 그 중 tks-contract는 고객의 계약 정보를 다루는 서비스이며, 다른 tks service들과 gRPC 기반으로 통신합니다. gRPC 호출을 위한 proto 파일은 [tks-proto](https://github.com/openinfradev/tks-proto)에서 확인할 수 있습니다.

## Quick Start

### Prerequisite
* docker 20.x 설치
* postgresql을 설치하고 database를 초기화합니다.
  ```
    docker run -p 5432:5432 --name postgres -e POSTGRES_PASSWORD=password -d postgres
    docker cp scripts/script.sql postgres:/script.sql
    docker exec -ti postgres psql -U postgres -a -f script.sql
  ``` 

### 서비스 구동 (For go developers)

```
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/tks-contract ./cmd/server/
$ bin/tks-contract -port 9110
```

### 서비스 구동 (For docker users)
```
$ docker pull sktcloud/tks-contract
$ docker run --name tks-contract -p 9110:9110 -d \
   sktcloud/tks-contract server -port 9110 
```

### gRPC API 호출 예제
```
import (

    "google.golang.org/grpc"
    pb "github.com/openinfradev/tks-proto/tks_pb"
    "google.golang.org/protobuf/encoding/protojson"
)

  func YOUR_FUNCTION(YOUR_PARAMS...) {
    server_addr = "tks-contract.taco-cat.xyz:9110"

    if len(args) == 0 {
        fmt.Println("Contract name must be specified.")
        fmt.Println("Usage: tksadmin contract create <CONTRACT NAME>")
        os.Exit(1)
    }

    var conn *grpc.ClientConn
    conn, err := grpc.Dial(server_addr, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Could not connect to server: %s", err)
    }
    defer conn.Close()

    client := pb.NewContractServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    //TODO: this quota is currently dummy and it'll updated later
    quota := &pb.ContractQuota{}
    quota.Cpu = 1200
    quota.Memory = 1200
    quota.Block = 1200
    quota.BlockSsd = 0
    quota.Fs = 0
    quota.FsSsd = 0

    data := pb.CreateContractRequest{
      ContractorName: args[0]
      Quota: quota
      AvailableServices: []string{"LMA", "SERVICE_MESH"}
      CspName: "test"
    }

    m := protojson.MarshalOptions{
        Indent:        "  ",
        UseProtoNames: true,
    }
    jsonBytes, _ := m.Marshal(&data)
    fmt.Println("Proto Json data...")
    fmt.Println(string(jsonBytes))

    r, err := client.CreateContract(ctx, &data)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Success: The request to create contract ", args[0], " was accepted.")
		}
	}

```
