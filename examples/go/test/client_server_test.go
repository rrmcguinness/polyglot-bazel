package test

import (
	"testing"

	grpc2 "example/grpc"
	"example/pkg/server"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var srv server.AuditServer
var clnt grpc2.AuditClient

func startServer(t *testing.T) {
	srv, err := server.NewAuditServer()
	assert.NotNil(t, srv)
	assert.Nil(t, err)
}

func stopServer(t *testing.T) {
	srv.Stop()
}

func getClient(t *testing.T) {
	var opts = make([]grpc.DialOption, 0)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(srv.GetConnectionString(), opts...)
	if err != nil {
		stopServer(t)
	}
	clnt = grpc2.NewAuditClient(conn)
}

func TestClientServer(t *testing.T) {
	startServer(t)
	getClient(t)
	assert.NotNil(t, clnt)

	stopServer(t)
}
