package models

import (
	"github.com/satori/go.uuid"
	"time"
)

// Event struct is the main internal representation of an event
type Event struct {
	EventID   uuid.UUID  `json:"event_id" db:"id"`
	UserName  string     `json:"user_name" db:"user_name"`
	EventName string     `json:"event_name" db:"title"`
	Note      string     `json:"note" db:"note"`
	StartTime *time.Time `json:"start_time" db:"start_time"`
	EndTime   *time.Time `json:"end_time" db:"end_time"`
}
