package state

import "github.com/libp2p/go-libp2p-core/peer"

type RoutingTable [][]peer.AddrInfo

func (r RoutingTable) Route(self, target peer.ID) *peer.AddrInfo {
	p := commonPrefix(self, target)

	if p >= len(r) {
		// @todo
	}

	//row := r[p]

	return nil
}

func commonPrefix(self, target peer.ID) int {
	return 0
}