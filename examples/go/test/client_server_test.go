package test

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"strconv"
	"testing"
	"time"

	"examples/go/pb"
	"examples/go/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestClientServer(t *testing.T) {
	assert.NotNil(t, Client)

	ctx := context.Background()

	messages := make([]*pb.Event, 0)
	for i := 0; i < 10; i++ {
		baggage := make(map[string]string)
		baggage["TEST_COUNT"] = strconv.Itoa(i)
		record := model.NewEvent("test", "go_test_client", "SYSTEM", baggage)
		messages = append(messages, record)
	}
	sent, err := Client.Put(ctx, nil, messages...)

	assert.Equal(t, 10, len(sent))
	if err != nil {
		log.Fatalf("failed to close stream : %v\n", err)
	}

	// Read the latest events
	req := &pb.DateRangeRequest{
		RequestId: model.NewRandomUUID(),
		Start:     timestamppb.Now(),
		End:       timestamppb.New(timestamppb.Now().AsTime().Add(24 * time.Hour)),
	}

	read, err := Client.FindByDateRange(ctx, nil, req)
	if err != nil {
		log.Fatalf("failed to read from bq: %v\n", err)
	}
	assert.NotNil(t, read)
	assert.True(t, len(read) > 0)
}
