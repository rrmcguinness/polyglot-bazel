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

// Package client simplifies the construction and use of the gRPC AuditClient.
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

type AuditClient struct {
	auditClient rpc.AuditClient
}

// NewAuditClient creates an audit client with default properties.
func NewAuditClient(conn *grpc.ClientConn) (*AuditClient, error) {
	out := &AuditClient{
		auditClient: rpc.NewAuditClient(conn),
	}
	return out, nil
}

// Ensures that the call options are not nil, and if they are returns
// an empty array.
func verifyCallOptions(callOptions []grpc.CallOption) []grpc.CallOption {
	if callOptions == nil {
		return make([]grpc.CallOption, 0)
	}
	return callOptions
}

// AuditCreateClient exposes access to the audit create client.
func (a *AuditClient) AuditCreateClient(
	ctx context.Context,
	callOptions []grpc.CallOption) (rpc.Audit_CreateClient, error) {
	callOptions = verifyCallOptions(callOptions)
	return a.auditClient.Create(ctx, callOptions...)
}

func Collect(createClient rpc.Audit_CreateClient) []*pb.AuditResponse {
	out := make([]*pb.AuditResponse, 0)
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

// Create adds one or more Audit Records to service and backing store.
func (a *AuditClient) Create(
	ctx context.Context,
	callOptions []grpc.CallOption,
	messages ...*pb.AuditRecord) ([]*pb.AuditResponse, error) {

	callOptions = verifyCallOptions(callOptions)
	createClient, err := a.AuditCreateClient(ctx, callOptions)

	for _, out := range messages {
		operation := func() (*pb.AuditRecord, error) {
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
	out := Collect(createClient)
	return out, err
}

func (a *AuditClient) Verify(
	ctx context.Context,
	callOptions []grpc.CallOption,
	responses ...*pb.AuditResponse) error {

	verifyClient, err := a.auditClient.Verify(ctx, verifyCallOptions(callOptions)...)
	if err != nil {
		return err
	}

	for _, response := range responses {
		operation := func() (*pb.AuditResponse, error) {
			if err = verifyClient.Send(response); err != nil {
				return response, err
			}
			return response, nil
		}
		err, r := backoff.RetryWithData(operation, backoff.NewExponentialBackOff())
		if err != nil {
			log.Printf("failed to get response: %v", r)
			// TODO - Add additional code here to fail over to a file for replay.
		}
	}
	return verifyClient.CloseSend()
}
