package state

import "github.com/libp2p/go-libp2p-core/peer"

type RoutingTable struct {
	rows [][]peer.AddrInfo
}
