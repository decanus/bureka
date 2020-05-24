package dht

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/decanus/bureka/pb"
)

type HandlerFunc func(ctx context.Context, from peer.ID, message pb.Message) *pb.Message

func (n *Node) handler(t pb.Message_Type) HandlerFunc {
	switch t {
	case pb.Message_MESSAGE:
		return n.onMessage
	case pb.Message_NODE_JOIN:
		return n.onNodeJoin
	case pb.Message_NODE_ANNOUNCE:
		return n.onNodeAnnounce
	case pb.Message_NODE_EXIT:
		return n.onNodeExit
	case pb.Message_HEARTBEAT:
		return n.onHeartbeat
	case pb.Message_REPAIR_REQUEST:
		return n.onRepairRequest
	case pb.Message_STATE_REQUEST:
		return n.onStateRequest
	case pb.Message_STATE_RESPONSE:
		return n.onStateRequest
	}

	return nil
}

func (n *Node) onMessage(ctx context.Context, _ peer.ID, message pb.Message) *pb.Message {
	err := n.Send(ctx, message)
	if err != nil {
		// @todo
	}

	return nil
}

func (n *Node) onNodeJoin(ctx context.Context, from peer.ID, message pb.Message) *pb.Message {
	// @TODO THIS IS QUESTIONABLE CAUSE IT MAY BE HANDLED THROUGH ANOTHER PATH ALREADY
	return nil
}

func (n *Node) onNodeAnnounce(ctx context.Context, from peer.ID, message pb.Message) *pb.Message {
	return nil
}

func (n *Node) onNodeExit(ctx context.Context, from peer.ID, message pb.Message) *pb.Message {
	err := n.remove(peer.ID(message.Sender))
	if err != nil {
		// @todo
	}

	return nil
}

func (n *Node) onHeartbeat(_ context.Context, _ peer.ID, message pb.Message) *pb.Message {
	for _, app := range n.applications {
		app.Heartbeat(peer.ID(message.Sender))
	}

	return nil
}

func (n *Node) onRepairRequest(ctx context.Context, from peer.ID, message pb.Message) *pb.Message {
	return nil
}

func (n *Node) onStateRequest(ctx context.Context, from peer.ID, message pb.Message) *pb.Message {
	return nil
}

func (n *Node) onStateResponse(ctx context.Context, from peer.ID, message pb.Message) *pb.Message {
	return nil
}
