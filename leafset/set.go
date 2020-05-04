package leafset

import (
	"bytes"
	"github.com/libp2p/go-libp2p-core/peer"
	k "github.com/libp2p/go-libp2p-kbucket"
	"sort"
)

// Set represents a Set of nodes
type Set []*k.PeerInfo

// Closest returns the closest peer to a specific ID.
func (s Set) Closest(id peer.ID) *k.PeerInfo {
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
func (s Set) Upsert(peer *k.PeerInfo) Set {
	i := s.search(peer.Id)
	if s[i].Id == peer.Id {
		s[i].LastSuccessfulOutboundQueryAt = peer.LastSuccessfulOutboundQueryAt
		s[i].LastUsefulAt = peer.LastUsefulAt
		return s
	}

	if i >= len(s) {
		s = append(s, nil)
		copy(s[i+1:], s[i:])
		s[i] = peer

		// @todo check its not too long
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
		if p.Id == id {
			return i
		}
	}

	return -1
}

func (s Set) search(id peer.ID) int {
	byteid, _ := id.MarshalBinary()

	return sort.Search(len(s), func(i int) bool {
		cmp, _ := (s)[i].Id.MarshalBinary()
		return bytes.Compare(byteid, cmp) >= 0
	})
}