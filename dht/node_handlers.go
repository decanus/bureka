package dht

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/decanus/bureka/pb"
)

type HandlerFunc func (ctx context.Context, from peer.ID, message pb.Message) error

func (n *Node) onMessage(ctx context.Context, from peer.ID, message pb.Message) {
	err := n.Send(message)
	if err != nil {
		// @todo
	}
}

func (n *Node) onNodeJoin(ctx context.Context, from peer.ID, message pb.Message) {
	// @TODO THIS IS QUESTIONABLE CAUSE IT MAY BE HANDLED THROUGH ANOTHER PATH ALREADY
}

func (n *Node) onNodeAnnounce(ctx context.Context, from peer.ID, message pb.Message) {

}

func (n *Node) onNodeExit(ctx context.Context, from peer.ID, message pb.Message) {

}

func (n *Node) onHeartbeat(ctx context.Context, from peer.ID, message pb.Message) {

}

func (n *Node) onRepairRequest(ctx context.Context, from peer.ID, message pb.Message) {

}

func (n *Node) onStateRequest(ctx context.Context, from peer.ID, message pb.Message) {

}

func (n *Node) onStateResponse(ctx context.Context, from peer.ID, message pb.Message) {

}
