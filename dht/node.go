package dht

import (
	"bufio"
	"context"
	"sync"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/routing"

	"github.com/decanus/bureka/state"
)

var logger = logging.Logger("dht")

const pastry protocol.ID = "/pastry/1.0/proto"

// Application represents a pastry application
type Application interface {
	Deliver(msg []byte)
	Forward(msg []byte, target peer.ID) bool
	Heartbeat(id peer.ID)
}

// Node is a pastry node.
type Node struct {
	sync.RWMutex

	LeafSet         state.LeafSet
	NeighborhoodSet state.Set
	RoutingTable    state.RoutingTable

	host host.Host

	applications []Application
}

// Guarantee that we implement interfaces.
var _ routing.PeerRouting = (*Node)(nil)

func New(ctx context.Context, host host.Host) (*Node, *Node) {
	return &Node{
		LeafSet:         state.NewLeafSet(host.ID()),
		NeighborhoodSet: make(state.Set, 0),
		applications:    make([]Application, 0),
		host:            host,
	}, nil
}

// AddApplication adds an application as a message receiver.
func (n *Node) AddApplication(app Application) {
	n.Lock()
	defer n.Unlock()

	n.applications = append(n.applications, app)
}

// Send sends a message to the target or the next closest peer.
func (n *Node) Send(ctx context.Context, msg []byte, key peer.ID) error {
	if key == n.host.ID() {
		n.deliver(msg) // @todo we may need to do this for more than just message types, like when the routing table is updated.
		return nil
	}

	target := n.route(key)
	if target.ID == "" {
		// no target to be found, delivering to self
		return nil
	}

	forward := n.forward(msg, target.ID)
	if !forward {
		return nil
	}

	err := n.send(ctx, msg, target.ID)
	if err != nil {
		return err
	}

	return nil
}

// ID returns a nodes ID, mainly for testing purposes.
func (n *Node) ID() peer.ID  {
	return n.host.ID()
}

func (n *Node) FindPeer(ctx context.Context, id peer.ID) (peer.AddrInfo, error) {
	if err := id.Validate(); err != nil {
		return peer.AddrInfo{}, err
	}

	logger.Debug("finding peer", "peer", id)

	local := n.route(id)
	if local.ID != "" {
		return local, nil
	}

	return peer.AddrInfo{}, nil
}

// @todo probably want to return error if not found
func (n *Node) route(to peer.ID) peer.AddrInfo {
	if n.LeafSet.IsInRange(to) {
		addr := n.LeafSet.Closest(to)
		if addr != nil {
			return *addr
		}
	}

	// @todo this is flimsy but will fix later
	addr := n.RoutingTable.Route(n.host.ID(), to)
	if addr != nil {
		return *addr
	}

	return peer.AddrInfo{}
}

// deliver sends the message to all connected applications.
func (n *Node) deliver(msg []byte) {
	n.RLock()
	defer n.RUnlock()

	for _, app := range n.applications {
		app.Deliver(msg)
	}
}

// forward asks all applications whether a message should be forwarded to a peer or not.
func (n *Node) forward(msg []byte, target peer.ID) bool {
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

func (n *Node) send(ctx context.Context, msg []byte, target peer.ID) error {
	s, err := n.host.NewStream(ctx, target, pastry)
	if err != nil {
		return err
	}

	bufw := bufio.NewWriter(s)

	// @todo probably needs this: https://github.com/libp2p/go-libp2p-pubsub/blob/5bbe37191afbb25a953e7931bf1a2ce18fbbb8f3/comm.go#L116

	_, err = bufw.Write(msg)
	if err != nil {
		return err
	}

	return nil
}
