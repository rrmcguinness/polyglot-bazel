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

package model

import (
	"examples/go/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewEvent(
	action string,
	context string,
	principal string,
	vars map[string]string) *pb.Event {

	return &pb.Event{
		Id:               NewRandomUUID(),
		Created:          timestamppb.Now(),
		Action:           action,
		Context:          context,
		Principal:        principal,
		ContextVariables: vars,
	}
}

// NewEventRecord is the constructor and requires the txId as it handles more than a single event.
func NewEventRecord(txId string) *pb.EventRecord {
	return &pb.EventRecord{
		TxId:     txId,
		Observed: timestamppb.Now(),
		Events:   make([]*pb.Event, 0),
	}
}
