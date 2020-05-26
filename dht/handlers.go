package dht

import (
	"context"

	"github.com/decanus/bureka/pb"
	"github.com/decanus/bureka/state"
)

type HandlerFunc func(ctx context.Context, from state.Peer, message *pb.Message) *pb.Message

func (n *Node) handler(t pb.Message_Type) HandlerFunc {
	return nil
}
