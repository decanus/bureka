package node

import (
	"context"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"

	"github.com/decanus/bureka/dht"
	"github.com/decanus/bureka/pb"
)

var logger = logging.Logger("dht")

// Node is a simple implementation that bridges libp2p IO to the pastry DHT state machine.
type Node struct {
	ctx context.Context

	dht    *dht.DHT
	host   host.Host
	writer *Writer
}

// Guarantee that we implement interfaces.
var _ routing.PeerRouting = (*Node)(nil)

func New(ctx context.Context, d *dht.DHT, h host.Host, w *Writer) (*Node, error) {
	return &Node{
		ctx: ctx,
		dht: d,
		host: h,
		writer: w,
	}, nil
}

func (n *Node) FindPeer(ctx context.Context, id peer.ID) (peer.AddrInfo, error) {
	if err := id.Validate(); err != nil {
		return peer.AddrInfo{}, err
	}

	logger.Debug("finding peer", "peer", id)

	b, _ := id.MarshalBinary()
	p := n.dht.Find(b)
	if p == nil {
		return peer.AddrInfo{}, nil // @todo error
	}

	id, err := peer.IDFromBytes(p)
	if err != nil {
		return peer.AddrInfo{}, err
	}

	return n.host.Peerstore().PeerInfo(id), nil
}

func (n *Node) Send(ctx context.Context, msg pb.Message) error {
	return n.dht.Send(ctx, msg)
}
