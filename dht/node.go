package dht

import (
	"context"
	"sync"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"

	"github.com/decanus/bureka/state"
)

var logger = logging.Logger("dht")

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

func New(ctx context.Context, host host.Host) *Node {
	return &Node{
		LeafSet:         state.NewLeafSet(host.ID()),
		NeighborhoodSet: make(state.Set, 0),
		applications:    make([]Application, 0),
	}
}

// AddApplication adds an application as a message receiver.
func (n *Node) AddApplication(app Application) {
	n.Lock()
	defer n.Unlock()

	n.applications = append(n.applications, app)
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
