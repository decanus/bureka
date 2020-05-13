package pastry

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/decanus/pastry/pb"
)

type HandlerFunc func (ctx context.Context, from peer.ID, message pb.Message) error

func (p *Pastry) onMessage(ctx context.Context, from peer.ID, message pb.Message) {
	err := p.Send(message)
	if err != nil {
		// @todo
	}
}

func (p *Pastry) onNodeJoin(ctx context.Context, from peer.ID, message pb.Message) {
	// @TODO THIS IS QUESTIONABLE CAUSE IT MAY BE HANDLED THROUGH ANOTHER PATH ALREADY
}

func (p *Pastry) onNodeAnnounce(ctx context.Context, from peer.ID, message pb.Message) {

}

func (p *Pastry) onNodeExit(ctx context.Context, from peer.ID, message pb.Message) {

}

func (p *Pastry) onHeartbeat(ctx context.Context, from peer.ID, message pb.Message) {

}

func (p *Pastry) onRepairRequest(ctx context.Context, from peer.ID, message pb.Message) {

}

func (p *Pastry) onStateRequest(ctx context.Context, pastry Pastry, from peer.ID, message pb.Message) {

}

func (p *Pastry) onStateResponse(ctx context.Context, pastry Pastry, from peer.ID, message pb.Message) {

}
