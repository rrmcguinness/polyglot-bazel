package test

import (
	"context"
	"log"
	"testing"

	"examples/go/pb"
	"examples/go/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestClientServer(t *testing.T) {
	assert.NotNil(t, Client)

	ctx := context.Background()

	messages := make([]*pb.AuditRecord, 0)
	for i := 1; i < 5; i++ {
		baggage := make(map[string]string)
		record := model.NewAuditRecord("test", "test_client", "SYSTEM", baggage)
		messages = append(messages, record)
	}
	sent, err := Client.Create(ctx, nil, messages...)

	assert.Equal(t, 4, len(sent))
	if err != nil {
		log.Fatalf("failed to close stream : %v\n", err)
	}
	log.Printf("Sent: %v\n", sent)
}
