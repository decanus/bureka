package dht

import (
	"context"

	peer "github.com/libp2p/go-libp2p-peer"

	"github.com/decanus/bureka/pb"
)

type HandlerFunc func(ctx context.Context, from peer.ID, message *pb.Message) *pb.Message

func (n *Node) handler(t pb.Message_Type) HandlerFunc {
	return nil
}
