package pastry

import (
	"github.com/libp2p/go-eventbus"
	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/multiformats/go-multiaddr"
)

func (p *Pastry) subscribe() {
	defer p.host.Network().StopNotify(p)

	evts := []interface{}{
		// register for event bus notifications of when peers successfully complete identification in order to update
		// the routing table
		new(event.EvtPeerIdentificationCompleted),

		// register for event bus protocol ID changes in order to update the routing table
		new(event.EvtPeerProtocolsUpdated),

		// register for event bus notifications for when our local address/addresses change so we can
		// advertise those to the network
		new(event.EvtLocalAddressesUpdated),
	}

	subs, err := p.host.EventBus().Subscribe(evts, eventbus.BufSize(256))

	for {
		e, ok := <-subs.Out()
		if !ok {
			return
		}

		switch event := e.(type) {
		case event.EvtPeerIdentificationCompleted:

		}
	}

}

func (p *Pastry) Listen(network network.Network, multiaddr multiaddr.Multiaddr)      {}
func (p *Pastry) ListenClose(network network.Network, multiaddr multiaddr.Multiaddr) {}
func (p *Pastry) Connected(network network.Network, conn network.Conn)               {}
func (p *Pastry) Disconnected(network network.Network, conn network.Conn)            {}
func (p *Pastry) OpenedStream(network network.Network, stream network.Stream)        {}
func (p *Pastry) ClosedStream(network network.Network, stream network.Stream)        {}
