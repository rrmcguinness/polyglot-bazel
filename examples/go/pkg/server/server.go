package server

import (
	"flag"
	"fmt"
	"log"
	"net"

	rpc "examples/go/grpc"
	"examples/go/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	host     = flag.String("host", "localhost", "The host interface to listen on")
	port     = flag.Int("port", 50080, "The server port")
)

func NewAuditServer() (*AuditGRPCServer, error) {
	out := &AuditGRPCServer{
		opts:       make([]grpc.ServerOption, 0),
		running:    false,
		tlsEnabled: *tls,
		certFile:   *certFile,
		keyFile:    *keyFile,
		Host:       *host,
		Port:       int32(*port),
		Log:        log.Default(),
	}

	out.PrepareGRPCOpts()
	out.server = grpc.NewServer(out.opts...)

	// Register all services with GRPC
	rpc.RegisterAuditServer(out.server, service.NewAuditService())

	return out, nil
}

type AuditGRPCServer struct {
	opts       []grpc.ServerOption
	running    bool
	tlsEnabled bool
	certFile   string
	keyFile    string
	Host       string
	Port       int32
	server     *grpc.Server
	quit       chan bool
	Log        *log.Logger
}

func (srv *AuditGRPCServer) ConnectionString() string {
	return fmt.Sprintf("%s:%d", srv.Host, srv.Port)
}

func (srv *AuditGRPCServer) fatalErrorCheck(msg string, err error) bool {
	if err != nil {
		srv.Log.Fatalf(msg, err)
		return false
	}
	return true
}

func (srv *AuditGRPCServer) PrepareGRPCOpts() {
	if *tls {
		if *certFile != "" && *keyFile != "" {
			creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
			if ok := srv.fatalErrorCheck("failed to generate credentials", err); ok {
				srv.opts = append(srv.opts, grpc.Creds(creds))
			}
		} else {
			// Anonymous credentials
		}
	}
}

// Start starts the server listening on the registered port.
func (srv *AuditGRPCServer) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", srv.Host, srv.Port))
	if ok := srv.fatalErrorCheck("failed to listen on port", err); ok {
		go func() {
			err = srv.server.Serve(lis)
			if ok := srv.fatalErrorCheck("failed to listen", err); ok {
				srv.running = true
				log.Printf("server started: %v", srv.ConnectionString())
			}
		}()
	}
}

// Stop gracefully shuts down the server
func (srv *AuditGRPCServer) Stop() {
	log.Printf("stopping srv: %s", srv.ConnectionString())
	srv.server.GracefulStop()
}
