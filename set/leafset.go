package set

import (
	"bytes"

	"github.com/libp2p/go-libp2p-core/peer"
)

// LeafSet contains the sets of numerically closer and farther from the node.
type LeafSet struct {
	key             peer.ID
	smaller, larger Set
}

// Closest returns the closest PeerInfo.
func (l LeafSet) Closest(id peer.ID) *peer.AddrInfo {
	byteid, _ := id.MarshalBinary()
	k, _ := l.key.MarshalBinary()
	if bytes.Compare(byteid, k) < 0 {
		return l.smaller.Closest(id)
	}

	return l.larger.Closest(id)
}

// Insert inserts a peer in the LeafSet.
func (l LeafSet) Insert(peer *peer.AddrInfo) {
	byteid, _ := peer.ID.MarshalBinary()
	k, _ := l.key.MarshalBinary()
	if bytes.Compare(byteid, k) < 0 {
		l.smaller = l.smaller.Insert(peer)
		return
	}

	l.larger = l.larger.Insert(peer)
}

// Remove removes a peer from the LeafSet.
func (l LeafSet) Remove(id peer.ID) bool {
	byteid, _ := id.MarshalBinary()
	k, _ := l.key.MarshalBinary()
	if bytes.Compare(byteid, k) < 0 {
		smaller, ok := l.smaller.Remove(id)
		l.smaller = smaller
		return ok
	}

	larger, ok := l.larger.Remove(id)
	l.larger = larger
	return ok
}

// @TODO not sure if this works lol, cause of left and right ordering.

// Min returns the farthest key to the smaller side.
func (l LeafSet) Min() peer.ID {
	return l.smaller[0].ID
}

// Max returns the farthest key to the larger side.
// @todo this may need to be l.larger[len(l.larger)-1].ID
func (l LeafSet) Max() peer.ID {
	return l.larger[0].ID
}
