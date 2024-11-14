package entity

import "github.com/google/uuid"

type ID struct {
	uuid.UUID
}

func NewID() ID {
	return ID{uuid.New()}
}

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID{id}, err
}
