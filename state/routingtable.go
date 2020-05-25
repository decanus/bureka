package state

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

// RoutingTable contains nodes organized by their distance to a peer.
type RoutingTable []Set

// Route returns the node closest to the target.
func (r RoutingTable) Route(self, target peer.ID) peer.ID {
	p := commonPrefix(self, target)

	if p >= len(r) {
		// @todo error handling
		return ""
	}

	return r[p].Closest(target)
}

// Insert adds a node to the RoutingTable.
func (r RoutingTable) Insert(self, id peer.ID) RoutingTable {
	nr := r
	p := commonPrefix(self, id)
	if p > len(r) {
		nr = r.grow(p)
	}

	nr[p] = nr[p].Insert(id)

	return nr
}

func (r RoutingTable) grow(n int) RoutingTable {
	nr := r
	if n > len(r) {
		appends := len(r) - n
		for i := 0; i <= appends; i++ {
			nr = append(r, make(Set, 0))
		}
	}

	return nr
}

func commonPrefix(self, target peer.ID) int {
	s, _ := self.MarshalBinary()
	t, _ := target.MarshalBinary()

	for i, v := range s {
		if v == t[i] {
			continue
		}

		return i
	}

	return len(s)
}
