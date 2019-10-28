package interfaces

import (
	"context"
	"github.com/omerkaya1/watcher/internal/models"
	"time"
)

// EventStorageProcessor is an interface to communicate with the DB
type EventStorageProcessor interface {
	// GetEventsByDate returns a slice of events that were created by the specified user
	GetUpcomingEvents(context.Context) ([]models.Event, error)
	// GetEventsForSpecifiedDate returns a slice of events that occur on a specified date
	GetEventsForSpecifiedDate(context.Context, time.Time) ([]models.Event, error)
}

// MessageQueueProcessor .
type MessageQueueProcessor interface {
	// ProduceMessages .
	ProduceMessages() error
}
