package repository

import (
	"github.com/google/uuid"
)

func ParseUUID(value string) (uuid.UUID, error) {
	return uuid.Parse(value)
}
