package pastry

import (
	"bytes"
	"github.com/libp2p/go-libp2p-core/peer"
	k "github.com/libp2p/go-libp2p-kbucket"
	"sort"
)

// Set represents a Set of nodes
type Set []*k.PeerInfo

// LeafSet contains the sets of numerically closer and farther from the node.
type LeafSet struct {
	smaller, larger Set
}

// Closest returns the closest peer to a specific ID.
func (s Set) Closest(id peer.ID) *k.PeerInfo {
	if len(s) == 0 {
		return nil
	}

	byteid, _ := id.MarshalBinary()

	i := sort.Search(len(s), func(i int) bool {
		cmp, _ := (s)[i].Id.MarshalBinary()
		return bytes.Compare(byteid, cmp) >= 0
	})

	if i >= len(s) {
		i = len(s) - 1
	}

	return (s)[i]
}

// Upsert either adds a peer to the Set or updates the peer if it already exists.
func (s Set) Upsert(peer *k.PeerInfo) Set {
	i := s.IndexOf(peer.Id)
	if i > -1 {
		return s
	}

	// @todo insertation
	return s
}

// Remove removes a peer with a given id.
func (s Set) Remove(id peer.ID) (Set, bool) {
	i := s.IndexOf(id)
	if i == -1 {
		return s, false
	}

	copy(s[i:], s[i+1:])
	s[len(s)-1] = nil
	return s[:len(s)-1], true
}

// IndexOf returns the index of the given peer id.
func (s Set) IndexOf(id peer.ID) int {
	for i, p := range s {
		if p.Id == id {
			return i
		}
	}

	return -1
}