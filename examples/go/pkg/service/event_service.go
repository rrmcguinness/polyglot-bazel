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
	"context"
	"errors"
	"examples/go/grpc"
	"examples/go/pb"
	"examples/go/pkg/model"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
)

var InvalidEventError = errors.New("invalid event, it requires an id and created attribute")

type EventService struct {
	grpc.UnimplementedEventsServer
}

func (svr *EventService) Get(ctx context.Context, request *pb.IdRequest) (*pb.EventRecord, error) {
	return nil, nil
}

func (svr *EventService) Put(stream grpc.Events_PutServer) error {
	// Collect events from stream
	record := &pb.EventRecord{
		TxId:     model.NewRandomUUID(),
		Observed: timestamppb.Now(),
		Events:   make([]*pb.Event, 0),
	}
	var err error

	for {
		event, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		if model.ValidateEvent(event) {
			err = InvalidEventError
			break
		}
		record.Events = append(record.Events, event)
		// Since this is an event, respond with an ack containing the transaction id
		// this may be used for accounting purposes such that a client send
		response := &pb.EventResponse{AuditRecordId: event.Id, TransactionId: record.TxId}
		err = stream.Send(response)
		if err != nil {
			break
		}
	}

	// Once closed, write to BigQuery
	if err == nil {
		err = insertEventRecord(record)
	}
	return err
}

func (svr *EventService) Verify(stream grpc.Events_VerifyServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return io.EOF
		}
		if err != nil {
			return err
		}
		// Check to see if the value is in BigQuery
		received := verify(in)
		if err := stream.Send(&pb.EventVerifyResponse{Verified: received}); err != nil {
			return err
		}
	}
}

func (svr *EventService) FindByDateRange(request *pb.DateRangeRequest, stream grpc.Events_FindByDateRangeServer) error {
	//if request != nil && request.End.AsTime().After(request.Start.AsTime()) {
	out, err := findByDateRange(request)
	if err != nil {
		return err
	}

	for _, r := range out {
		err = stream.Send(r)
		if err != nil {
			break
		}
	}

	return err
}
