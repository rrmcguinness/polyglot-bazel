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
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ValidateEvent(event *pb.Event) bool {
	if event == nil {
		return false
	}
	// Parse the UUID, ensure it's valid
	_, err := uuid.FromBytes([]byte(event.Id))
	if err != nil {
		return false
	}
	return event.Created != nil
}

type ServiceEvent struct {
	Id               string      `json:"id" bigquery:"id"`
	Created          time.Time   `json:"created" bigquery:"created"`
	Action           string      `json:"action" bigquery:"action"`
	Context          string      `json:"context" bigquery:"context"`
	Principal        string      `json:"principal" bigquery:"principal"`
	ContextVariables []*KeyValue `json:"context_variables" bigquery:"context_variables"`
}

func NewServiceEvent(event *pb.Event) *ServiceEvent {
	out := &ServiceEvent{
		Id:               event.Id,
		Created:          event.Created.AsTime(),
		Action:           event.Action,
		Context:          event.Context,
		Principal:        event.Principal,
		ContextVariables: make([]*KeyValue, 0),
	}
	for k, v := range event.ContextVariables {
		if len(k) > 0 && len(v) > 0 {
			out.ContextVariables = append(out.ContextVariables, &KeyValue{Key: k, Value: v})
		}
	}
	return out
}

type KeyValue struct {
	Key   string `json:"k" bigquery:"k"`
	Value string `json:"v" bigquery:"v"`
}

type ServiceEventRecord struct {
	TxId     string          `json:"tx_id" bigquery:"tx_id"`
	Observed time.Time       `json:"observed" bigquery:"observed"`
	Events   []*ServiceEvent `json:"events" bigquery:"events"`
}

func NewServiceEventRecord(event *pb.EventRecord) *ServiceEventRecord {
	out := &ServiceEventRecord{
		TxId:     event.TxId,
		Observed: event.Observed.AsTime(),
		Events:   make([]*ServiceEvent, 0),
	}
	for _, e := range event.Events {
		out.Events = append(out.Events, NewServiceEvent(e))
	}
	return out
}

func EventRecordFromServiceEventRecord(record *ServiceEventRecord) *pb.EventRecord {
	out := NewEventRecord(record.TxId)
	out.Observed = timestamppb.New(record.Observed)
	for _, e := range record.Events {
		mapVals := make(map[string]string)
		for _, k := range e.ContextVariables {
			mapVals[k.Key] = k.Value
		}
		out.Events = append(out.Events, &pb.Event{
			Id:               e.Id,
			Created:          timestamppb.New(e.Created),
			Action:           e.Action,
			Context:          e.Context,
			Principal:        e.Principal,
			ContextVariables: mapVals,
		})
	}
	return out
}
