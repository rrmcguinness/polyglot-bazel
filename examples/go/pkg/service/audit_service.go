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
	"log"

	"example/grpc"
	"example/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewAuditService() *AuditService {
	return &AuditService{
		UnimplementedAuditServer: grpc.UnimplementedAuditServer{},
		records:                  make(map[string]*pb.AuditRecord),
	}
}

type AuditService struct {
	grpc.UnimplementedAuditServer
	records map[string]*pb.AuditRecord
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
		svc.records[txId.String()] = audit

		// Since this is an audit, respond with an ack containing the transaction id
		// this may be used for accounting purposes such that a client send
		// 100 audit events, and the service verified those events.
		stream.Send(&pb.AuditResponse{
			AuditRecordId: audit.Id,
			TransactionId: txId.String(),
		})
	}
}
