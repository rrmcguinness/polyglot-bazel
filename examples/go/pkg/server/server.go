// Package server is a default server for GRPC functions. Each GRPC service
// must be registered in the "NewAuditServer" method.
package server

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

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

// NewAuditServer is the creator method for the GRPC Server
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
		quit:       make(chan bool),
	}

	out.prepareGRPCOpts()
	out.server = grpc.NewServer(out.opts...)

	// Register all services with GRPC
	rpc.RegisterAuditServer(out.server, service.NewAuditService())

	return out, nil
}

// AuditGRPCServer is the default server for audit functions.
type AuditGRPCServer struct {
	opts       []grpc.ServerOption
	running    bool
	tlsEnabled bool
	quit       chan bool
	certFile   string
	keyFile    string
	Host       string
	Port       int32
	mu         sync.Mutex
	server     *grpc.Server
	Log        *log.Logger
}

// Safely sets the server state for other observers
func (srv *AuditGRPCServer) setRunning(val bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.running = val
}

// IsRunning is a convenience method for determining if the server is running.
func (srv *AuditGRPCServer) IsRunning() bool {
	return srv.running
}

// ConnectionString A convenience method for getting host and port.
func (srv *AuditGRPCServer) ConnectionString() string {
	return fmt.Sprintf("%s:%d", srv.Host, srv.Port)
}

// Handle errors gracefully
func (srv *AuditGRPCServer) fatalErrorCheck(msg string, err error) bool {
	if err != nil {
		srv.Log.Fatalf(msg, err)
		return false
	}
	return true
}

// Sets up the credentials from the given cert and key files.
func (srv *AuditGRPCServer) prepareGRPCOpts() {
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

func runServer(srv *AuditGRPCServer, lis net.Listener) {
	// This will block until Close is called.
	err := srv.server.Serve(lis)
	if err != nil {
		srv.fatalErrorCheck("failed to listen", err)
		return
	}
	return
}

func stopServer(srv *AuditGRPCServer) {
	srv.Log.Printf("stopping server: %s", srv.ConnectionString())
	srv.server.GracefulStop()
	srv.setRunning(false)
}

// Start starts the server in the background listening on the registered port.
func (srv *AuditGRPCServer) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", srv.Host, srv.Port))
	if ok := srv.fatalErrorCheck("failed to listen on port", err); ok {
		go func() {
			for {
				select {
				case <-srv.quit:
					stopServer(srv)
					return
				default:
					if !srv.running {
						go runServer(srv, lis)
					}
				}
			}
		}()
	}
}

// Stop gracefully shuts down the server
func (srv *AuditGRPCServer) Stop() {
	srv.quit <- true
}
