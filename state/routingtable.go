package state

import (
	"bytes"
	"sort"

	"github.com/libp2p/go-libp2p-core/peer"
)

type RoutingTable [][]peer.ID

// Route returns the node closest to the target.
func (r RoutingTable) Route(self, target peer.ID) peer.ID {
	p := commonPrefix(self, target)

	if p >= len(r) {
		// @todo error handling
		return ""
	}

	row := r[p]

	b, _ := target.MarshalBinary()

	// @todo this may be wrong
	// see: https://github.com/secondbit/wendy/blob/e4601da9fbf96fd1f6e81a18e58db10b57bce3ff/nodeid.go#L214
	if row[b[p]] != "" {
		return row[b[p]]
	}

	// @todo find node closer numerically

	return ""
}

// Insert adds a node to the RoutingTable.
func (r RoutingTable) Insert(self, id peer.ID) RoutingTable {
	p := commonPrefix(self, id)
	if p > len(r) {
		r = r.grow(p)
	}

	row := r[p]

	byteid, _ := id.MarshalBinary()

	i := sort.Search(len(row), func(i int) bool {
		cmp, _ := (row)[i].MarshalBinary()
		return bytes.Compare(byteid, cmp) >= 0
	})

	if i < len(row) && row[i] == id || i >= Length {
		return r
	}

	nr := append(row, "")
	copy(nr[i+1:], nr[i:])
	nr[i] = id

	r[p] = nr

	return r
}

func (r RoutingTable) grow(n int) RoutingTable {
	if n > len(r) {
		appends := len(r) - n
		for i := 0; i <= appends; i++ {
			r = append(r, make([]peer.ID, 0))
		}
	}

	return r
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
