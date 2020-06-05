package dht

import (
	"bytes"
	"context"
	"sync"

	"github.com/decanus/bureka/dht/state"
	"github.com/decanus/bureka/pb"
)

// Transport is responsible for sending messages.
// This represents a call back function that can be implemented on network IO.
type Transport interface {
	Send(ctx context.Context, target state.Peer, msg *pb.Message) error
}

// ApplicationID represents a unique identifier for the application.
type ApplicationID string

// Application represents a pastry application
type Application interface {
	Deliver(msg *pb.Message)
	Forward(msg *pb.Message, target state.Peer) bool
	Heartbeat(id state.Peer)
}

// DHT represents a pastry DHT in its most basic form as a state machine.
type DHT struct {
	sync.RWMutex

	ID state.Peer

	LeafSet         state.LeafSet
	NeighborhoodSet state.Set
	RoutingTable    state.RoutingTable

	applications map[ApplicationID]Application

	transport Transport
}

// New returns a new DHT.
func New(id state.Peer, transport Transport) *DHT {
	return &DHT{
		ID:              id,
		LeafSet:         state.NewLeafSet(id),
		NeighborhoodSet: make(state.Set, 0),
		RoutingTable:    make(state.RoutingTable, 0),
		applications:    make(map[ApplicationID]Application),
		transport:       transport,
	}
}

// AddApplication adds an application as a message receiver.
func (d *DHT) AddApplication(aid ApplicationID, app Application) {
	d.Lock()
	defer d.Unlock()

	d.applications[aid] = app
}

// RemoveApplication removes an application from the set.
func (d *DHT) RemoveApplication(aid ApplicationID) {
	d.Lock()
	defer d.Unlock()

	delete(d.applications, aid)
}

// Send a message to the target peer or closest available peer.
func (d *DHT) Send(ctx context.Context, msg *pb.Message) error {
	key := msg.Key

	if bytes.Equal(key, d.ID) {
		d.deliver(msg) // @todo we may need to do this for more than just message types, like when the routing table is updated.
		return nil
	}

	target := d.Find(key)
	if target == nil {
		d.deliver(msg)
		return nil
	}

	forward := d.forward(msg, target)
	if !forward {
		return nil
	}

	err := d.transport.Send(ctx, target, msg)
	if err != nil {
		return err
	}

	return nil
}

// Find returns the closest known peer to a given target or the target itself.
func (d *DHT) Find(target state.Peer) state.Peer {
	d.RLock()
	defer d.RUnlock()

	if d.LeafSet.IsInRange(target) {
		id := d.LeafSet.Closest(target)
		if id != nil {
			return id
		}
	}

	// @todo this is flimsy but will fix later
	id := d.RoutingTable.Route(d.ID, target)
	if id != nil {
		return id
	}

	return nil
}

// AddPeer adds a newly found peer to the dht.
func (d *DHT) AddPeer(id state.Peer) {
	d.Lock()
	defer d.Unlock()

	// @todo probably need to think about max length for neighborhoodset
	d.NeighborhoodSet = d.NeighborhoodSet.Insert(id)
	d.RoutingTable = d.RoutingTable.Insert(d.ID, id)
	d.LeafSet.Insert(id)
}

// RemovePeer removes a peer from the dht.
func (d *DHT) RemovePeer(id state.Peer) {
	d.Lock()
	defer d.Unlock()

	ns, _ := d.NeighborhoodSet.Remove(id)
	d.NeighborhoodSet = ns

	d.RoutingTable = d.RoutingTable.Remove(d.ID, id)
	d.LeafSet.Remove(id)
}

func (d *DHT) Heartbeat(id state.Peer) {
	d.RLock()
	defer d.RUnlock()

	for _, app := range d.applications {
		app.Heartbeat(id)
	}
}

// MapNeighbors iterates over the NeighborhoodSet and calls the process for every peer.
func (d *DHT) MapNeighbors(process func(peer state.Peer)) {
	d.RLock()
	defer d.RUnlock()

	for _, p := range d.NeighborhoodSet {
		process(p)
	}
}

func (d *DHT) MapRoutingTable(process func(peer state.Peer)) {
	d.RLock()
	defer d.RUnlock()

	for _, r := range d.RoutingTable {
		for _, p := range r {
			process(p)
		}
	}
}

// deliver sends the message to all connected applications.
func (d *DHT) deliver(msg *pb.Message) {
	d.RLock()
	defer d.RUnlock()

	for _, app := range d.applications {
		app.Deliver(msg)
	}
}

// forward asks all applications whether a message should be forwarded to a peer or not.
func (d *DHT) forward(msg *pb.Message, target state.Peer) bool {
	d.RLock()
	defer d.RUnlock()

	// @todo need to run over this logic
	forward := true
	for _, app := range d.applications {
		f := app.Forward(msg, target)
		if forward {
			forward = f
		}
	}

	return forward
}
