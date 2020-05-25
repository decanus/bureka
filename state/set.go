package state

import (
	"bytes"
	"sort"

	"github.com/libp2p/go-libp2p-core/peer"
)

var Length int = 10

// Set represents a Set of nodes
type Set []peer.ID

// Closest returns the closest peer to a specific ID.
func (s Set) Closest(id peer.ID) peer.ID {
	if len(s) == 0 {
		return ""
	}

	i := s.search(id)

	if i >= len(s) {
		i = len(s) - 1
	}

	return (s)[i]
}

// Insert adds a peer to the Set.
func (s Set) Insert(peer peer.ID) Set {
	i := s.search(peer)

	if i < len(s) && s[i] == peer || i >= Length {
		return s
	}

	ns := append(s, "")
	copy(ns[i+1:], ns[i:])
	ns[i] = peer

	return ns
}

// Remove removes a peer with a given id.
func (s Set) Remove(id peer.ID) (Set, bool) {
	i := s.IndexOf(id)
	if i == -1 {
		return s, false
	}

	copy(s[i:], s[i+1:])
	s[len(s)-1] = ""
	return s[:len(s)-1], true
}

// IndexOf returns the index of the given peer id.
func (s Set) IndexOf(id peer.ID) int {
	for i, p := range s {
		if p == id {
			return i
		}
	}

	return -1
}

func (s Set) search(id peer.ID) int {
	byteid, _ := id.MarshalBinary()

	return sort.Search(len(s), func(i int) bool {
		cmp, _ := (s)[i].MarshalBinary()
		return bytes.Compare(byteid, cmp) >= 0
	})
}
