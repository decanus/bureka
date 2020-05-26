package state

import (
	"bytes"
	"sort"
)

var Length int = 10

type Peer []byte

// Set represents a Set of nodes
type Set []Peer

// Closest returns the closest peer to a specific ID.
func (s Set) Closest(id Peer) Peer {
	if len(s) == 0 {
		return nil
	}

	i := s.search(id)

	if i >= len(s) {
		i = len(s) - 1
	}

	return (s)[i]
}

// Insert adds a peer to the Set.
func (s Set) Insert(peer Peer) Set {
	i := s.search(peer)

	if i < len(s) && bytes.Equal(s[i], peer) || i >= Length {
		return s
	}

	ns := append(s, nil)
	copy(ns[i+1:], ns[i:])
	ns[i] = peer

	return ns
}

// Remove removes a peer with a given id.
func (s Set) Remove(id Peer) (Set, bool) {
	i := s.IndexOf(id)
	if i == -1 {
		return s, false
	}

	copy(s[i:], s[i+1:])
	s[len(s)-1] = nil
	return s[:len(s)-1], true
}

// IndexOf returns the index of the given peer id.
func (s Set) IndexOf(id Peer) int {
	for i, p := range s {
		if bytes.Equal(p, id) {
			return i
		}
	}

	return -1
}

func (s Set) search(id Peer) int {

	return sort.Search(len(s), func(i int) bool {
		return bytes.Compare(id, (s)[i]) >= 0
	})
}
