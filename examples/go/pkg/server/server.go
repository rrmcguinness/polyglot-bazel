package server

import (
	"flag"
	"fmt"
	"log"
	"net"

	"example/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	host     = flag.String("host", "localhost", "The host interface to listen on")
	port     = flag.Int("port", 50051, "The server port")
)

func NewAuditServer() (*AuditServer, error) {
	server := &AuditServer{
		opts:       make([]grpc.ServerOption, 0),
		running:    false,
		tlsEnabled: *tls,
		certFile:   *certFile,
		keyFile:    *keyFile,
		Host:       *host,
		Port:       int32(*port),
		quit:       make(chan bool),
		Log:        log.Default(),
	}

	return server, nil
}

type AuditServer struct {
	opts         []grpc.ServerOption
	running      bool
	tlsEnabled   bool
	certFile     string
	keyFile      string
	Host         string
	Port         int32
	server       *grpc.Server
	quit         chan bool
	auditService *service.AuditService
	Log          *log.Logger
}

func (server *AuditServer) GetConnectionString() string {
	return server.Host + ":" + string(server.Port)
}

func (server *AuditServer) fatalErrorCheck(msg string, err error) bool {
	if err != nil {
		server.Log.Fatalf(msg, err)
		return false
	}
	return true
}

func (server *AuditServer) PrepareGRPCOpts() {
	if *tls {
		if *certFile != "" && *keyFile != "" {
			creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
			if ok := server.fatalErrorCheck("failed to generate credentials", err); ok {
				server.opts = append(server.opts, grpc.Creds(creds))
			}
		} else {
			// Anonymous credentials
		}
	}
}

// Start starts the server listening on the registered port.
func (server *AuditServer) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port))
	if ok := server.fatalErrorCheck("failed to listen on port", err); ok {
		go func() {
			for {
				select {
				case <-server.quit:
					log.Printf("Stopping Server: %s:%d", server.Host, server.Port)
					server.server.GracefulStop()
					server.running = false
					return
				default:
					err = server.server.Serve(lis)
					if ok := server.fatalErrorCheck("failed to listen", err); ok {
						server.running = true
					}
				}
			}
		}()
	}
}

// Stop gracefully shuts down the server
func (server *AuditServer) Stop() {
	server.quit <- true
}
