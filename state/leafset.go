package state

import (
	"bytes"
)

// LeafSet contains the sets of numerically closer and farther from the node.
type LeafSet struct {
	key             Peer
	smaller, larger Set
}

func NewLeafSet(key Peer) LeafSet {
	return LeafSet{
		key:     key,
		smaller: make(Set, 0),
		larger:  make(Set, 0),
	}
}

// Insert inserts a peer in the LeafSet.
func (l *LeafSet) Insert(peer Peer) {
	if bytes.Compare(peer, l.key) < 0 {
		l.smaller = l.smaller.Insert(peer)
		return
	}

	l.larger = l.larger.Insert(peer)
}

// Remove removes a peer from the LeafSet.
func (l *LeafSet) Remove(id Peer) bool {
	if bytes.Compare(id, l.key) < 0 {
		smaller, ok := l.smaller.Remove(id)
		l.smaller = smaller
		return ok
	}

	larger, ok := l.larger.Remove(id)
	l.larger = larger
	return ok
}

// Closest returns the closest PeerInfo.
func (l LeafSet) Closest(id Peer) Peer {
	if bytes.Compare(id, l.key) < 0 {
		return l.smaller.Closest(id)
	}

	return l.larger.Closest(id)
}

// Min returns the farthest key to the smaller side.
func (l LeafSet) Min() Peer {
	if len(l.smaller) == 0 {
		return nil
	}

	return l.smaller[len(l.smaller)-1]
}

// Max returns the farthest key to the larger side.
func (l LeafSet) Max() Peer {
	if len(l.larger) == 0 {
		return nil
	}

	return l.larger[0]
}

// IsInRange returns whether an id is between
// the Min and Max IDs in the LeafSet.
func (l LeafSet) IsInRange(id Peer) bool {
	min := l.Min()
	max := l.Max()

	if min != nil && bytes.Compare(id, min) < 1 {
		return false
	}

	if max != nil && bytes.Compare(id, max) > 0 {
		return false
	}

	return true
}