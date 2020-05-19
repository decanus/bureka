package dht

import (
	"context"
	"sync"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/routing"

	"github.com/decanus/bureka/pb"
	"github.com/decanus/bureka/state"
)

var logger = logging.Logger("dht")
const pastry = protocol.ID("/pastry/1.0/proto")

// Application represents a pastry application
type Application interface {
	Deliver(message *pb.Message)
	Forward(message *pb.Message, target peer.ID) bool
	Heartbeat(id peer.ID)
}

// Node implements the main logic of the DHT.
// This includes managing the LeafSet, NeighborhoodSet, RoutingTable
// as well as dealing with messages.
type Node struct {
	sync.RWMutex

	ctx context.Context

	LeafSet         state.LeafSet
	NeighborhoodSet state.Set
	RoutingTable    state.RoutingTable

	host host.Host

	applications []Application
}

// Guarantee that we implement interfaces.
var _ routing.PeerRouting = (*Node)(nil)

func New(ctx context.Context, host host.Host) *Node {
	n := &Node{
		ctx:             ctx,
		LeafSet:         state.NewLeafSet(host.ID()),
		NeighborhoodSet: make(state.Set, 0),
		host:            host,
		applications:    make([]Application, 0),
	}

	n.host.SetStreamHandler(pastry, n.streamHandler)

	return n
}

// AddApplication adds an application as a message receiver.
func (n *Node) AddApplication(app Application) {
	n.Lock()
	defer n.Unlock()

	n.applications = append(n.applications, app)
}

// Send delivers a message to the next closest target.
func (n *Node) Send(msg *pb.Message) error {
	key := peer.ID(msg.Key)

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

	err := n.send(msg, target.ID)
	if err != nil {
		return err
	}

	return nil
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

// deliver sends the message to all connected applications.
func (n *Node) deliver(msg *pb.Message) {
	n.RLock()
	defer n.RUnlock()

	for _, app := range n.applications {
		app.Deliver(msg)
	}
}

// forward asks all applications whether a message should be forwarded to a peer or not.
func (n *Node) forward(msg *pb.Message, target peer.ID) bool {
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

func (n *Node) send(msg *pb.Message, target peer.ID) error {
	// @todo
	return nil
}

func (n *Node) remove(peer peer.ID) error {
	n.LeafSet.Remove(peer)
	n.NeighborhoodSet.Remove(peer)

	// @todo from routing table.

	return nil
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
