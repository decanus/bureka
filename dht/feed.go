package dht

import (
	"sync"

	"github.com/decanus/bureka/dht/state"
	"github.com/decanus/bureka/pb"
)

// Packet is used to pass messages and their targets in feeds.
type Packet struct {
	Target  state.Peer
	Message *pb.Message
}

// Subscription represents a feed channel.
type Subscription chan<- Packet

// Feed implements one-to-many subscriptions where the carrier of events is a channel.
//
// It was inspired by: https://github.com/prysmaticlabs/prysm/blob/546196a6fa8174d68f5c92071ff4bc6edd1ce3d0/shared/event/feed.go
type Feed struct {
	sync.Mutex

	subscribers []Subscription
}

// NewFeed returns a new feed.
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
	f.Lock()
	defer f.Unlock()

	// @todo is this good enough for now?
	for _, sub := range f.subscribers {
		sub <- value
	}
}
