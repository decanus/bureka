package state

import (
	"bytes"
	"sort"
)

// Peer type represents a peer id.
type Peer []byte

// @todo think about changing set struct to have a max max field
// Set represents a Set of nodes
type Set struct {
	peers []Peer
	max   int
}

func NewSet(length int) Set {
	return Set{
		peers: make([]Peer, 0),
		max: length,
	}
}

// Closest returns the closest peer to a specific ID.
func (s Set) Closest(id Peer) Peer {
	if len(s.peers) == 0 {
		return nil
	}

	i := s.insertAt(id)

	if i >= len(s.peers) {
		i = len(s.peers) - 1
	}

	return s.peers[i]
}

// Insert adds a peer to the Set.
func (s Set) Insert(peer Peer) Set {
	i := s.insertAt(peer)

	if i < len(s.peers) && bytes.Equal(s.peers[i], peer) {
		return s
	}

	// @todo, what we could do here is init a Set at a certain max and never append.
	if len(s.peers) < s.max {
		s.peers = append(s.peers, nil)
	}

	copy(s.peers[i+1:], s.peers[i:])
	s.peers[i] = peer

	return s
}

// Remove removes a peer with a given id.
func (s Set) Remove(id Peer) (Set, bool) {
	i := s.IndexOf(id)
	if i == -1 {
		return s, false
	}

	copy(s.peers[i:], s.peers[i+1:])
	s.peers[len(s.peers)-1] = nil
	s.peers = s.peers[:len(s.peers)-1]

	return s, true
}

// IndexOf returns the index of the given peer id.
func (s Set) IndexOf(id Peer) int {
	for i, p := range s.peers {
		if bytes.Equal(p, id) {
			return i
		}
	}

	return -1
}

// Get returns the value at the index in the set.
func (s Set) Get(index int) Peer {
	return s.peers[index]
}

// Length returns the max of the set
func (s Set) Length() int {
	return len(s.peers)
}

func (s Set) Map(process func(peer Peer)) {
	for _, p := range s.peers {
		process(p)
	}
}

func (s Set) insertAt(id Peer) int {
	return sort.Search(len(s.peers), func(i int) bool {
		return bytes.Compare(id, s.peers[i]) >= 0
	})
}
