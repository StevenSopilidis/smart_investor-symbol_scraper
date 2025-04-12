package domain

import (
	"github.com/google/uuid"
)

type Symbol struct {
	ID       uuid.UUID
	Ticker   string
	Exchange string
	Active   bool
}
