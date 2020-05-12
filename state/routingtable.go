package state

import "github.com/libp2p/go-libp2p-core/peer"

type RoutingTable [][]peer.AddrInfo

func (r RoutingTable) Route(self, target peer.ID) *peer.AddrInfo {
	// @TODO
	return nil
}
