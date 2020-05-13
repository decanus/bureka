// Package dht implements a pastry node.
//
// The implementation is inspired by https://github.com/libp2p/go-libp2p-kad-dht,
// as well as various Node implementations including https://github.com/secondbit/wendy.
package dht

import (
	"context"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"

	"github.com/decanus/pastry"
	"github.com/decanus/pastry/state"
)

var logger = logging.Logger("dht")

type Node struct {
	LeafSet         state.LeafSet
	NeighborhoodSet state.Set
	RoutingTable    state.RoutingTable

	host host.Host

	deliverHandler bureka.DeliverHandler
	forwardHandler bureka.ForwardHandler
}

// Guarantee that we implement interfaces.
var _ routing.PeerRouting = (*Node)(nil)

func New(ctx context.Context, host host.Host) *Node {
	return &Node{
		LeafSet:         state.NewLeafSet(host.ID()),
		NeighborhoodSet: make(state.Set, 0),
	}
}

func (n *Node) FindPeer(ctx context.Context, id peer.ID) (peer.AddrInfo, error) {
	if err := id.Validate(); err != nil {
		return peer.AddrInfo{}, err
	}

	logger.Debug("finding peer", "peer", id)

	local := n.route(id)
	if local.ID != "" {
		return local, nil
	}

	return peer.AddrInfo{}, nil
}

// @todo probably want to return error if not found
func (n *Node) route(to peer.ID) peer.AddrInfo {
	if n.LeafSet.IsInRange(to) {
		addr := n.LeafSet.Closest(to)
		if addr != nil {
			return *addr
		}
	}

	// @todo this is flimsy but will fix later
	addr := n.RoutingTable.Route(n.host.ID(), to)
	if addr != nil {
		return *addr
	}

	return peer.AddrInfo{}
}
