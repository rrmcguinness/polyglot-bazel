/*
 * Copyright 2023 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"io"
	"sync"

	"examples/go/grpc"
	"examples/go/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewAuditService() *auditServer {
	return &auditServer{
		records: make(map[string]*pb.AuditRecord),
	}
}

type auditServer struct {
	grpc.UnimplementedAuditServer
	mu      sync.Mutex
	records map[string]*pb.AuditRecord
}

func (svr *auditServer) Verify(stream grpc.Audit_VerifyServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return io.EOF
		}
		if err != nil {
			return err
		}
		received := svr.records[in.TransactionId]
		if err := stream.Send(&pb.AuditVerificationResponse{Verified: received != nil}); err != nil {
			return err
		}
	}
}

func (svr *auditServer) Create(stream grpc.Audit_CreateServer) error {
	for {
		audit, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		//log.Printf("%v", audit)

		persist := pb.AuditRecord{
			Id:               audit.Id,
			Created:          audit.Created,
			Action:           audit.Action,
			Context:          audit.Context,
			Principal:        audit.Principal,
			ContextVariables: make(map[string]string),
		}

		// Copy the variables
		for k, v := range audit.ContextVariables {
			persist.ContextVariables[k] = v
		}

		if len(persist.Id) == 0 {
			id, _ := uuid.NewRandom()
			persist.Id = id.String()
		}

		if persist.Created == nil {
			persist.Created = timestamppb.Now()
		}

		txId, _ := uuid.NewRandom()
		svr.mu.Lock()
		svr.records[txId.String()] = audit
		svr.mu.Unlock()

		// Since this is an audit, respond with an ack containing the transaction id
		// this may be used for accounting purposes such that a client send
		response := &pb.AuditResponse{
			AuditRecordId: audit.Id, TransactionId: txId.String()}

		stream.Send(response)
	}
}
