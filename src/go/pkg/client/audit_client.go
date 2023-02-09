package client

import (
	"io"
	"log"

	"example/grpc"
	"example/pb"
)

type Callable interface {
	Callable()
	GetWaitChannel() chan struct{}
	GetCommChannel() chan *pb.AuditResponse
}

type AuditRecordResponseCallback struct {
	waitc  chan struct{}
	stream grpc.Audit_CreateClient
	comm   chan *pb.AuditResponse
}

func (audit AuditRecordResponseCallback) GetCommChannel() chan *pb.AuditResponse {
	return audit.comm
}

func (audit AuditRecordResponseCallback) GetWaitChannel() chan struct{} {
	return audit.waitc
}

func (audit AuditRecordResponseCallback) Callback() {
	for {
		in, err := audit.stream.Recv()
		if err == io.EOF {
			close(audit.waitc)
			return
		}
		if err != nil {
			log.Default().Printf("error while receiving: %v", err)
		}
		audit.comm <- in
	}
}

func NewAuditRecordResponseCallback(stream grpc.Audit_CreateClient) *AuditRecordResponseCallback {
	return &AuditRecordResponseCallback{
		waitc:  make(chan struct{}),
		stream: stream,
		comm:   make(chan *pb.AuditResponse),
	}
}

type AuditClient struct {
	grpc.AuditClient
	handler AuditRecordResponseCallback
}
