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

// Package client simplifies the construction and use of the gRPC EventClient.
// intended to supply hold the client state and manage the connections.
package client

import (
	"context"
	"io"
	"log"

	rpc "examples/go/grpc"
	"examples/go/pb"
	"github.com/cenkalti/backoff/v4"
	"google.golang.org/grpc"
)

type EventClient struct {
	eventClient rpc.EventsClient
}

// NewEventClient creates an audit client with default properties.
func NewEventClient(conn *grpc.ClientConn) *EventClient {
	out := &EventClient{
		eventClient: rpc.NewEventsClient(conn),
	}
	return out
}

// Ensures that the call options are not nil, and if they are returns
// an empty array.
func verifyCallOptions(callOptions []grpc.CallOption) []grpc.CallOption {
	if callOptions == nil {
		return make([]grpc.CallOption, 0)
	}
	return callOptions
}

func collect(createClient rpc.Events_PutClient) []*pb.EventResponse {
	out := make([]*pb.EventResponse, 0)
	waitc := make(chan struct{})
	go func() {
		for {
			m, e := createClient.Recv()
			if e == io.EOF {
				close(waitc)
				break
			}
			if e != nil {
				close(waitc)
				log.Fatalf("failed reading the client: %v\n", e)
			}
			out = append(out, m)
		}
	}()
	<-waitc
	return out
}

// Put adds one or more Event Records to service and backing store.
func (a *EventClient) Put(
	ctx context.Context,
	callOptions []grpc.CallOption,
	messages ...*pb.Event) ([]*pb.EventResponse, error) {

	createClient, err := a.eventClient.Put(ctx, verifyCallOptions(callOptions)...)

	for _, out := range messages {
		operation := func() (*pb.Event, error) {
			if err = createClient.Send(out); err != nil {
				return out, err
			}
			return out, nil
		}
		r, err := backoff.RetryWithData(operation, backoff.NewExponentialBackOff())
		if err != nil {
			log.Printf("failed to write record: %v with error: %v", r, err)
		}
	}
	err = createClient.CloseSend()
	out := collect(createClient)
	return out, err
}

func (a *EventClient) Verify(
	ctx context.Context,
	callOptions []grpc.CallOption,
	responses ...*pb.EventResponse) error {

	verifyClient, err := a.eventClient.Verify(ctx, verifyCallOptions(callOptions)...)
	if err != nil {
		return err
	}

	for _, response := range responses {
		operation := func() (*pb.EventResponse, error) {
			if err = verifyClient.Send(response); err != nil {
				return response, err
			}
			return response, nil
		}
		err, r := backoff.RetryWithData(operation, backoff.NewExponentialBackOff())
		if err != nil {
			log.Printf("failed to get response: %v", r)
		}
	}
	return verifyClient.CloseSend()
}

func (a *EventClient) FindByDateRange(
	ctx context.Context,
	opts []grpc.CallOption,
	in *pb.DateRangeRequest) ([]*pb.EventRecord, error) {

	resp, err := a.eventClient.FindByDateRange(ctx, in, verifyCallOptions(opts)...)

	if err != nil {
		return nil, err
	}

	out := make([]*pb.EventRecord, 0)

	for {
		r, err := resp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}
