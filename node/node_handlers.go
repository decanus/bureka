package node

import (
	"context"

	"github.com/gogo/protobuf/proto"

	"github.com/decanus/bureka/dht/state"
	"github.com/decanus/bureka/pb"
)

type handlerFunc func(ctx context.Context, message *pb.Message) (*pb.Message, error)

func (n *Node) handler(t pb.Message_Type) handlerFunc {
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
	case pb.Message_STATE_REQUEST:
		return n.onStateRequest
	case pb.Message_STATE_DATA:
		return n.onStateData
	}

	return nil
}

func (n *Node) onMessage(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	err := n.Send(ctx, message)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (n *Node) onNodeJoin(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	n.dht.AddPeer(message.Sender)

	resp, err := n.stateResponseMessage()
	if err != nil {
		return nil, err
	}

	resp.Key = message.Sender
	return resp, err
}

func (n *Node) onNodeAnnounce(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	n.dht.AddPeer(message.Sender)
	return nil, nil
}

func (n *Node) onNodeExit(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	n.dht.RemovePeer(message.Sender)
	return nil, nil
}

func (n *Node) onHeartbeat(_ context.Context, message *pb.Message) (*pb.Message, error) {
	n.dht.Heartbeat(message.Sender)
	return nil, nil
}

// @todo make sure this is what we want
func (n *Node) onStateRequest(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	resp, err := n.stateResponseMessage()
	if err != nil {
		return nil, err
	}

	resp.Key = message.Sender
	return resp, err
}

func (n *Node) onStateData(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	req := &pb.State{}
	err := proto.Unmarshal(message.Data, req)
	if err != nil {
		return nil, err
	}

	n.dht.ImportPeers(req.RoutingTable, req.Neighborhood, req.Leafset)

	return nil, nil
}

// stateResponseMessage returns a message with the current state tables.
func (n *Node) stateResponseMessage() (*pb.Message, error) {
	routing := make([][]byte, 0)
	n.dht.MapRoutingTable(func(peer state.Peer) {
		routing = append(routing, peer)
	})

	neighbor := make([][]byte, 0)
	n.dht.MapRoutingTable(func(peer state.Peer) {
		neighbor = append(neighbor, peer)
	})

	leafset := make([][]byte, 0)
	n.dht.MapLeafSet(func(peer state.Peer) {
		leafset = append(leafset, peer)
	})

	s := &pb.State{
		Neighborhood: neighbor,
		Leafset:      leafset,
		RoutingTable: routing,
	}

	d, err := proto.Marshal(s)
	if err != nil {
		return nil, err
	}

	return &pb.Message{
		Type:   pb.Message_STATE_DATA,
		Sender: n.dht.ID,
		Data:   d,
	}, nil
}
