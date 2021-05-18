# tks-contract

[![Go Report Card](https://goreportcard.com/badge/github.com/sktelecom/tks-contract?style=flat-square)](https://goreportcard.com/report/github.com/sktelecom/tks-contract)
[![Go Reference](https://pkg.go.dev/badge/github.com/sktelecom/tks-contract.svg)](https://pkg.go.dev/github.com/sktelecom/tks-contract)
[![Release](https://img.shields.io/github/release/sktelecom/tks-contract.svg?style=flat-square)](https://github.com/sktelecom/tks-contract/releases/latest)

Tks-contract is one of the tks services. It mainly deals with the contract information.  
This service communicates based on gRPC. You can refer to the proto files in [tks-proto](https://github.com/sktelecom/tks-proto).

## Quick Start

### Development environment
* Installed docker 20.x
* Running postgresql and Initilizing database.
  ```
    docker run -p 5432:5432 --name postgres -e POSTGRES_PASSWORD=password -d postgres
    docker cp scripts/script.sql postgres:/script.sql
    docker exec -ti postgres psql -U postgres -a -f script.sql
  ``` 
### For go developers

```
go install -v ./...
contract-server -port 9110
```
### For docker users
```
TAGS=$(curl --silent "https://api.github.com/repos/sktelecom/tks-contract/tags" | grep name | head -1 |cut -d '"' -f 4)
docker pull docker.pkg.github.com/sktelecom/tks-contract/tks-contract:$TAGS
docker run --name tks-contract -p 9110:9110 -d \
  docker.pkg.github.com/sktelecom/tks-contract/tks-contract:$TAGS \
  contract-server \
  # -port 9110
```

