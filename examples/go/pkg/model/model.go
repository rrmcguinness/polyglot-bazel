package model

import (
	"examples/go/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewAuditRecord(
	action string,
	context string,
	principal string,
	vars map[string]string) *pb.AuditRecord {

	id, _ := uuid.NewRandom()
	created := timestamppb.Now()

	return &pb.AuditRecord{
		Id:               id.String(),
		Created:          created,
		Action:           action,
		Context:          context,
		Principal:        principal,
		ContextVariables: vars,
	}
}
