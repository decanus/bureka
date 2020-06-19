package dht

import (
	"github.com/decanus/bureka/dht/state"
	"github.com/decanus/bureka/pb"
)

type Packet struct {
	Target  state.Peer
	Message *pb.Message
}

type Subscription chan<- Packet

type Feed struct {
	subscribers []Subscription
}

func NewFeed() *Feed {
	return &Feed{
		subscribers: make([]Subscription, 0),
	}
}

// Subscribe adds a channel to the feed.
func (f *Feed) Subscribe(channel Subscription) { // @todo think about returning a subscription like prysm
	f.subscribers = append(f.subscribers, channel)
}

// Send sends a payload to all the subscribers for the specific feed.
func (f *Feed) Send(value Packet) {
	// @todo is this good enough for now?
	for _, sub := range f.subscribers {
		sub <- value
	}
}
