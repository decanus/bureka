package leafset

import (
	"bytes"
	"github.com/libp2p/go-libp2p-core/peer"
	k "github.com/libp2p/go-libp2p-kbucket"
)

// LeafSet contains the sets of numerically closer and farther from the node.
type LeafSet struct {
	key []byte // @todo better type
	smaller, larger Set
}

// Closest returns the closest PeerInfo.
func (l LeafSet) Closest(id peer.ID) *k.PeerInfo {
	byteid, _ := id.MarshalBinary()
	if bytes.Compare(byteid, l.key) < 0 {
		return l.smaller.Closest(id)
	}

	return l.larger.Closest(id)
}

// Upsert either Insert or Updates a peer in the LeafSet.
func (l LeafSet) Upsert(peer *k.PeerInfo) {
	byteid, _ := peer.Id.MarshalBinary()
	if bytes.Compare(byteid, l.key) < 0 {
		l.smaller = l.smaller.Upsert(peer)
		return
	}

	l.larger = l.larger.Upsert(peer)
}

// Remove removes a peer from the LeafSet.
func (l LeafSet) Remove(id peer.ID) bool {
	byteid, _ := id.MarshalBinary()
	if bytes.Compare(byteid, l.key) < 0 {
		smaller, ok := l.smaller.Remove(id)
		l.smaller = smaller
		return ok
	}

	larger, ok := l.larger.Remove(id)
	l.larger = larger
	return ok
}
