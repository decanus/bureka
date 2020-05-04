package leafset

import (
	"bytes"
	"github.com/libp2p/go-libp2p-core/peer"
	"sort"
)

var SetLength int = 10

// Set represents a Set of nodes
type Set []*peer.AddrInfo

// Closest returns the closest peer to a specific ID.
func (s Set) Closest(id peer.ID) *peer.AddrInfo{
	if len(s) == 0 {
		return nil
	}

	i := s.search(id)

	if i >= len(s) {
		i = len(s) - 1
	}

	return (s)[i]
}

// Upsert either adds a peer to the Set or updates the peer if it already exists.
func (s Set) Upsert(peer *peer.AddrInfo) Set {
	i := s.search(peer.ID)
	if i >= SetLength {
		return s
	}

	if i >= len(s) {
		s = append(s, nil)
		copy(s[i+1:], s[i:])
		s[i] = peer

		if len(s) > SetLength {
			return s[:len(s)-1]
		}
	}

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
		if p.ID == id {
			return i
		}
	}

	return -1
}

func (s Set) search(id peer.ID) int {
	byteid, _ := id.MarshalBinary()

	return sort.Search(len(s), func(i int) bool {
		cmp, _ := (s)[i].ID.MarshalBinary()
		return bytes.Compare(byteid, cmp) >= 0
	})
}