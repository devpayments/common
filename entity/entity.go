package entity

import "github.com/google/uuid"

func NewNullUUID(uuidString string) *uuid.NullUUID {
	parsedUUID, err := uuid.Parse(uuidString)
	if err != nil {
		return nil
	}
	return &uuid.NullUUID{UUID: parsedUUID, Valid: true}
}

type Entity interface {
}

type Model[E Entity] interface {
	ToEntity() E
	FromEntity(entity E) interface{}
}
