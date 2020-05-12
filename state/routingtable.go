package state

import "github.com/libp2p/go-libp2p-core/peer"

type RoutingTable [][]*peer.AddrInfo

func (r RoutingTable) Route(self, target peer.ID) *peer.AddrInfo {
	p := commonPrefix(self, target)

	if p >= len(r) {
		// @todo error handling
		return nil
	}

	row := r[p]

	b, _ := target.MarshalBinary()

	// @todo this may be wrong
	// see: https://github.com/secondbit/wendy/blob/e4601da9fbf96fd1f6e81a18e58db10b57bce3ff/nodeid.go#L214
	if row[b[p]] != nil {
		return row[b[p]]
	}

	// @todo find node closer numerically

	return nil
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