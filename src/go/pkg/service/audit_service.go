package service

import (
	"io"
	"log"

	"example/grpc"
)

type AuditService struct {
	grpc.UnimplementedAuditServer
}

func (svc *AuditService) Create(stream grpc.Audit_CreateServer) (err error) {
	for {
		audit, err := stream.Recv()
		if err == io.EOF {
			return io.EOF
		}
		if err != nil {
			return err
		}
		log.Default().Printf("%v", audit)
	}
}
