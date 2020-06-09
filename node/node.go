package node

import (
	"bytes"
	"context"
	"sync"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"

	"github.com/decanus/bureka/dht"
	"github.com/decanus/bureka/dht/state"
	"github.com/decanus/bureka/node/internal"
	"github.com/decanus/bureka/pb"
)

var logger = logging.Logger("dht")

// ApplicationID represents a unique identifier for the application.
type ApplicationID string

// Application represents a pastry application
type Application interface {
	Deliver(msg *pb.Message)
	Forward(msg *pb.Message, target state.Peer) bool
	Heartbeat(id state.Peer)
}

// Node is a simple implementation that bridges libp2p IO to the pastry DHT state machine.
type Node struct {
	sync.RWMutex

	ctx context.Context

	dht    *dht.DHT
	host   host.Host
	writer *internal.Writer

	sub event.Subscription

	applications map[ApplicationID]Application
}

// Guarantee that we implement interfaces.
var _ routing.PeerRouting = (*Node)(nil)

// New returns a new Node.
func New(ctx context.Context, d *dht.DHT, h host.Host) (*Node, error) {
	n := &Node{
		ctx:    ctx,
		dht:    d,
		host:   h,
		writer: internal.NewWriter(),
	}

	s, err := n.subscribe()
	if err != nil {
		return nil, err
	}

	n.sub = s

	// adds the already known peers
	for _, p := range n.host.Network().Peers() {
		n.dht.AddPeer([]byte(p))
	}

	go n.poll(n.sub)

	return n, nil
}

// FindPeer finds the closest AddrInfo to the passed ID.
func (n *Node) FindPeer(ctx context.Context, id peer.ID) (peer.AddrInfo, error) {
	if err := id.Validate(); err != nil {
		return peer.AddrInfo{}, err
	}

	logger.Debug("finding peer", "peer", id)

	b := []byte(id)
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

func (n *Node) Send(ctx context.Context, msg *pb.Message) error {
	key := msg.Key

	if bytes.Equal(key, []byte(n.host.ID())) {
		n.deliver(msg) // @todo we may need to do this for more than just message types, like when the routing table is updated.
		return nil
	}

	target := n.dht.Find(key)
	if target == nil {
		n.deliver(msg)
		return nil
	}

	forward := n.forward(msg, target)
	if !forward {
		return nil
	}

	err := n.writer.Send(ctx, target, msg)
	if err != nil {
		return err
	}

	return nil
}

// deliver sends the message to all connected applications.
func (n *Node) deliver(msg *pb.Message) {
	n.RLock()
	defer n.RUnlock()

	for _, app := range n.applications {
		app.Deliver(msg)
	}
}

// forward asks all applications whether a message should be forwarded to a peer or not.
func (n *Node) forward(msg *pb.Message, target state.Peer) bool {
	n.RLock()
	defer n.RUnlock()

	// @todo need to run over this logic
	forward := true
	for _, app := range n.applications {
		f := app.Forward(msg, target)
		if forward {
			forward = f
		}
	}

	return forward
}
