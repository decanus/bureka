package pastry

import (
	"bytes"
	"context"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/decanus/pastry/set"
)

var logger = logging.Logger("dht")

type Pastry struct {
	LeafSet          set.LeafSet
	NeighbourhoodSet set.Set

	deliverHandler DeliverHandler
	forwardHandler ForwardHandler
}

func (p *Pastry) Route(ctx context.Context, to peer.ID) {
	if isInRange(to, p.LeafSet.Min(), p.LeafSet.Max()) {
		closest := p.LeafSet.Closest(to)
		// @todo route to closest
	} else {
		// @todo use routing table
	}
}

func (p *Pastry) FindPeer(ctx context.Context, id peer.ID) (peer.AddrInfo, error) {
	if err := id.Validate(); err != nil {
		return peer.AddrInfo{}, err
	}

	logger.Debug("finding peer", "peer", id)

	local := p.FindLocal(id)
	if local != nil {
		return *local, nil
	}

	// @todo should probably call route here.

	return peer.AddrInfo{}, nil
}

func (p *Pastry) FindLocal(id peer.ID) *peer.AddrInfo {
	// @todo should probably call route here.
	closest := p.LeafSet.Closest(id)
	if closest.ID == id {
		return closest
	}

	return nil
}

func (p *Pastry) route(ctx context.Context, to peer.ID) peer.AddrInfo {
	if isInRange(to, p.LeafSet.Min(), p.LeafSet.Max()) {
		addr := p.LeafSet.Closest(to)
		if addr != nil {
			return *addr
		}
	} else {
		// @todo use routing table
	}

	// @todo
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
