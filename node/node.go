package node

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"

	"github.com/decanus/bureka/dht"
	"github.com/decanus/bureka/dht/state"
	"github.com/decanus/bureka/pb"
)

type Node struct {
	dht  dht.DHT
	host host.Host
}

func New() (*Node, error) {
	return nil, nil
}

func (n *Node) Send(ctx context.Context, target state.Peer, msg pb.Message) error {
	panic("implement me")
}

