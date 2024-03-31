package vo

import "github.com/google/uuid"

type ID struct {
	value uuid.UUID
}

func (i ID) Value() uuid.UUID {
	return i.value
}
