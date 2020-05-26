package dht

import (
	"context"

	"github.com/decanus/bureka/dht/state"
	"github.com/decanus/bureka/pb"
)

type HandlerFunc func(ctx context.Context, from state.Peer, message *pb.Message) *pb.Message

func (n *Node) handler(t pb.Message_Type) HandlerFunc {
	return nil
}
