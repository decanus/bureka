package state

import "github.com/libp2p/go-libp2p-core/peer"

type RoutingTable [][]Peer

// Route returns the node closest to the target.
func (r RoutingTable) Route(self, target Peer) Peer {
	p := commonPrefix(self, target)

	if p >= len(r) {
		// @todo error handling
		return nil
	}

	row := r[p]

	// @todo this may be wrong
	// see: https://github.com/secondbit/wendy/blob/e4601da9fbf96fd1f6e81a18e58db10b57bce3ff/nodeid.go#L214
	if row[target[p]] != nil {
		return row[target[p]]
	}

	// @todo find node closer numerically

	return nil
}

func (r RoutingTable) Insert(self, id peer.ID) RoutingTable {
	// @todo
	return r
}

func commonPrefix(self, target Peer) int {
	for i, v := range self {
		if v == target[i] {
			continue
		}

		return i
	}

	return len(self)
}
