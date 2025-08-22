package events

// Event defines the interface for domain eventbus in DDD.
type Event interface {
	Type() string
}
