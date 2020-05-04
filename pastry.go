package pastry

import (
	"context"
	"github.com/decanus/pastry/leafset"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/peer"
)

var logger = logging.Logger("dht")

type Pastry struct {
	LeafSet leafset.LeafSet
}

func (p *Pastry) Route() {

}

func (p *Pastry) Deliver() {

}

func (p *Pastry) Forward() {

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

	return peer.AddrInfo{}, nil
}

func (p *Pastry) FindLocal(id peer.ID) *peer.AddrInfo {
	closest := p.LeafSet.Closest(id)
	if closest.Id == id {
		return closest
	}

	return nil
}