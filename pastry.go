// Package pastry implements a pastry node.
//
// The implementation is inspired by https://github.com/libp2p/go-libp2p-kad-dht,
// as well as various Pastry implementations including https://github.com/secondbit/wendy.
package pastry

import (
	"context"
	"sync"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/routing"

	"github.com/decanus/pastry/pb"
	"github.com/decanus/pastry/state"
)

var logger = logging.Logger("dht")
var proto = protocol.ID("/pastry/1.0/proto")

// Application represents a pastry application
type Application interface {
	Deliver(message pb.Message)
	Forward(message pb.Message, target peer.ID) bool
}

type Pastry struct {
	sync.RWMutex

	LeafSet         state.LeafSet
	NeighborhoodSet state.Set
	RoutingTable    state.RoutingTable

	host host.Host

	applications []Application
}

// Guarantee that we implement interfaces.
var _ routing.PeerRouting = (*Pastry)(nil)

func New(ctx context.Context, host host.Host) *Pastry {
	p := &Pastry{
		LeafSet:         state.NewLeafSet(host.ID()),
		NeighborhoodSet: make(state.Set, 0),
	}

	p.host.SetStreamHandler(proto, p.streamHandler)

	return p
}

func (p *Pastry) Send(msg pb.Message) error {
	key := peer.ID(msg.Key)

	if key == p.host.ID() {
		p.deliver(msg) // @todo we may need to do this for more than just message types, like when the routing table is updated.
		return nil
	}

	target := p.route(key)
	if target.ID == "" {
		// no target to be found, delivering to self
		return nil
	}

	forward := p.forward(msg, target.ID)
	if !forward {
		return nil
	}

	err := p.send(msg, target.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pastry) FindPeer(ctx context.Context, id peer.ID) (peer.AddrInfo, error) {
	if err := id.Validate(); err != nil {
		return peer.AddrInfo{}, err
	}

	logger.Debug("finding peer", "peer", id)

	local := p.route(id)
	if local.ID != "" {
		return local, nil
	}

	return peer.AddrInfo{}, nil
}

// deliver sends the message to all connected applications.
func (p *Pastry) deliver(msg pb.Message) {
	p.RLock()
	defer p.RUnlock()

	for _, app := range p.applications {
		app.Deliver(msg)
	}
}

// forward asks all applications whether a message should be forwarded to a peer or not.
func (p *Pastry) forward(msg pb.Message, target peer.ID) bool {
	p.RLock()
	defer p.RUnlock()

	// @todo need to run over this logic
	forward := true
	for _, app := range p.applications {
		f := app.Forward(msg, target)
		if forward {
			forward = f
		}
	}

	return forward
}

func (p *Pastry) send(msg pb.Message, target peer.ID) error {
	// @todo
	return nil
}

// @todo probably want to return error if not found
func (p *Pastry) route(to peer.ID) peer.AddrInfo {
	if p.LeafSet.IsInRange(to) {
		addr := p.LeafSet.Closest(to)
		if addr != nil {
			return *addr
		}
	}

	// @todo this is flimsy but will fix later
	addr := p.RoutingTable.Route(p.host.ID(), to)
	if addr != nil {
		return *addr
	}

	return peer.AddrInfo{}
}
