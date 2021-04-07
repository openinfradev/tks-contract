package contract

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	pb "github.com/openinfradev/tks-proto/pbgo"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Contract struct {
	ContractorName    string            `json:"contractorName"`
	ContractId        ContractId        `json:"contractId"`
	AvailableServices []string          `json:"availableServices,omitempty"`
	McOpsId           McOpsId           `json:"mcOpsId,omitempty"`
	LastUpdatedTs     *LastUpdatedTime  `json:"lastUpdatedTs"`
	Quota             *pb.ContractQuota `json:"quota"`
}

type ContractId string
type McOpsId uuid.UUID
type LastUpdatedTime struct {
	time.Time
}

// UnmarshalJSON parses the Unix time and stores the result in ts
func (ts *LastUpdatedTime) UnmarshalJSON(data []byte) error {
	unix, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	ts.Time = time.Unix(unix, 0)
	return nil
}

func (m McOpsId) String() string {
	return m.String()
}

func (l LastUpdatedTime) Timestamppb() *timestamppb.Timestamp {
	return timestamppb.New(l.Time)
}

func GenerateMcOpsId() McOpsId {
	return McOpsId(uuid.Must(uuid.NewRandom()))
}
