package errors

type WatcherError string

func (we WatcherError) Error() string {
	return string(we)
}

var (
	ErrBadConfigFile      = WatcherError("the correct configuration file was not specified")
	ErrBadDBConfiguration = WatcherError("malformed or uninitialised DB configuration")
	// Event conflict related errors
	ErrEventCollisionInInterval = WatcherError("event takes place within the time interval of another event")
	ErrEventDoesNotExist        = WatcherError("the requested event does not exist in the DB")
	ErrEventTimeViolation       = WatcherError("new events cannot be created in the past")
	ErrMalformedTimeObject      = WatcherError("invalid time string")
	// DB related errors
	ErrNoOpDBAction = WatcherError("no rows were affected by the action")
	// Other errors
	ErrBadQueueConfiguration = WatcherError("malformed or uninitialised message queue configuration")
)

const (
	ErrCMDPrefix = "command failure"
	ErrDBPrefix  = "db failure"
	ErrMQPrefix  = "message queue failure"
)
