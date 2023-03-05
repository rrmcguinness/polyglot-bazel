package test

import (
	"log"
	"os"
	"testing"

	"examples/go/pkg/client"
	"examples/go/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Server *server.EventsGRPCServer
var Client *client.EventClient
var conn *grpc.ClientConn

func serverSetup() {
	var err error
	Server, err = server.NewEventsServer()
	if err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

func clientSetup() {
	var err error
	conn, err = grpc.Dial(
		Server.ConnectionString(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to create client connection: %v", err)
	}

	Client = client.NewEventClient(conn)
}

func tearDown() {
	// Close the clients
	if conn != nil {
		err := conn.Close()
		if err != nil {
			log.Printf("Failed to close GRPC Connection to Server: %v", err)
		}
	}
	// Close the server
	if Server != nil {
		Server.Stop()
	}
}

func TestMain(t *testing.M) {
	serverSetup()
	Server.Start()
	clientSetup()
	sig := t.Run()
	tearDown()
	os.Exit(sig)
}
