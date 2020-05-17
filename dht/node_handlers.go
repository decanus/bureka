package dht

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/decanus/bureka/pb"
)

type HandlerFunc func (ctx context.Context, from peer.ID, message *pb.Message) *pb.Message

func (n *Node) onMessage(ctx context.Context, _ peer.ID, message *pb.Message) *pb.Message {
	err := n.Send(message)
	if err != nil {
		// @todo
	}

	return nil
}

func (n *Node) onNodeJoin(ctx context.Context, from peer.ID, message *pb.Message) *pb.Message {
	// @TODO THIS IS QUESTIONABLE CAUSE IT MAY BE HANDLED THROUGH ANOTHER PATH ALREADY
	return nil
}

func (n *Node) onNodeAnnounce(ctx context.Context, from peer.ID, message *pb.Message) *pb.Message {
	return nil
}

func (n *Node) onNodeExit(ctx context.Context, from peer.ID, message *pb.Message) *pb.Message {
	err := n.remove(peer.ID(message.Sender))
	if err != nil {
		// @todo
	}

	return nil
}

func (n *Node) onHeartbeat(ctx context.Context, _ peer.ID, message *pb.Message) *pb.Message {
	for _, app := range n.applications {
		app.Heartbeat(peer.ID(message.Sender))
	}

	return nil
}

func (n *Node) onRepairRequest(ctx context.Context, from peer.ID, message *pb.Message) *pb.Message {
	return nil
}

func (n *Node) onStateRequest(ctx context.Context, from peer.ID, message *pb.Message) *pb.Message {
	return nil
}

func (n *Node) onStateResponse(ctx context.Context, from peer.ID, message *pb.Message) *pb.Message {
	return nil
}
