package model

import (
	"github.com/google/uuid"
)

func NewRandomUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}
