package state

import "github.com/libp2p/go-libp2p-core/peer"

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
