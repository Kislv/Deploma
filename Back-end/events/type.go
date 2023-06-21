package events

import "context"

type Fetcher interface {
	Fetch(ctx context.Context, limit int) ([]Event, error)
}

type Processor interface {
	Process(ctx context.Context, e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}

type LesionParameters struct {
	Area float32
	Diameter float32
}

type ImageParameters struct {
	Extension string
	Height int32
	Width int32
	Data []byte
}