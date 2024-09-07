package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(event Event) error
}

type Type int

const (
	unknown Type = iota
	Message
)

type Event struct {
	Type Type
}
