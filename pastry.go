// Package pastry implements a pastry node.
//
// The implementation is inspired by [go-libp2p-kad-dht](https://github.com/libp2p/go-libp2p-kad-dht),
// as well as various Pastry implementations including [wendy](https://github.com/secondbit/wendy).
package pastry

import (
	"bytes"
	"context"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/decanus/pastry/state"
)

var logger = logging.Logger("dht")

type Pastry struct {
	LeafSet          state.LeafSet
	NeighbourhoodSet state.Set
	RoutingTable     state.RoutingTable

	host host.Host

	deliverHandler DeliverHandler
	forwardHandler ForwardHandler
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

func (p *Pastry) route(to peer.ID) peer.AddrInfo {
	if isInRange(to, p.LeafSet.Min(), p.LeafSet.Max()) {
		addr := p.LeafSet.Closest(to)
		if addr != nil {
			return *addr
		}
	}

	addr := p.RoutingTable.Route(p.host.ID(), to)
	if addr != nil {
		return *addr
	}

	return peer.AddrInfo{}
}

func isInRange(id, min, max peer.ID) bool {
	byteid, _ := id.MarshalBinary()
	bytemin, _ := min.MarshalBinary()
	if bytes.Compare(byteid, bytemin) >= 0 {
		return false
	}

	bytemax, _ := max.MarshalBinary()
	if bytes.Compare(byteid, bytemax) < 1 {
		return false
	}

	return true
}
