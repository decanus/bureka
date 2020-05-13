// Package pastry implements a pastry node.
//
// The implementation is inspired by https://github.com/libp2p/go-libp2p-kad-dht,
// as well as various Pastry implementations including https://github.com/secondbit/wendy.
package pastry

import (
	"context"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/routing"

	"github.com/decanus/pastry/state"
)

var logger = logging.Logger("dht")
var proto = protocol.ID("/ipfs-onion/1.0/proto")

type Pastry struct {
	LeafSet         state.LeafSet
	NeighborhoodSet state.Set
	RoutingTable    state.RoutingTable

	host host.Host

	deliverHandler DeliverHandler
	forwardHandler ForwardHandler
}

// Guarantee that we implement interfaces.
var _ routing.PeerRouting = (*Pastry)(nil)

func New(ctx context.Context, host host.Host) *Pastry {
	p := &Pastry{
		LeafSet:         state.NewLeafSet(host.ID()),
		NeighborhoodSet: make(state.Set, 0),
	}

	p.host.SetStreamHandler(proto, p.streamHandler)

	return p
}

func (p *Pastry) FindPeer(ctx context.Context, id peer.ID) (peer.AddrInfo, error) {
	if err := id.Validate(); err != nil {
		return peer.AddrInfo{}, err
	}

	logger.Debug("finding peer", "peer", id)

	local := p.route(id)
	if local.ID != "" {
		return local, nil
	}

	return peer.AddrInfo{}, nil
}

// @todo probably want to return error if not found
func (p *Pastry) route(to peer.ID) peer.AddrInfo {
	if p.LeafSet.IsInRange(to) {
		addr := p.LeafSet.Closest(to)
		if addr != nil {
			return *addr
		}
	}

	// @todo this is flimsy but will fix later
	addr := p.RoutingTable.Route(p.host.ID(), to)
	if addr != nil {
		return *addr
	}

	return peer.AddrInfo{}
}
