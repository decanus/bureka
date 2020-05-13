package pastry

import (
	"github.com/libp2p/go-eventbus"
	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

func (n *Node) subscribe() error {
	defer n.host.Network().StopNotify(n)

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

	subs, err := n.host.EventBus().Subscribe(evts, eventbus.BufSize(256))
	if err != nil {
		return err
	}

	for {
		e, ok := <-subs.Out()
		if !ok {
			return nil // @todo?
		}

		switch evt := e.(type) {
		case event.EvtPeerIdentificationCompleted:
			n.handlePeerIdentificationCompleted(evt.Peer)
		}
	}

}

func (n *Node) handlePeerIdentificationCompleted(id peer.ID) {
	info := n.host.Peerstore().PeerInfo(id)
	if info.ID == "" {
		return // @todo
	}

	n.discovered(&info)
}

func (n *Node) Listen(network network.Network, multiaddr multiaddr.Multiaddr)      {}
func (n *Node) ListenClose(network network.Network, multiaddr multiaddr.Multiaddr) {}
func (n *Node) Connected(network network.Network, conn network.Conn)               {}
func (n *Node) Disconnected(network network.Network, conn network.Conn)            {}
func (n *Node) OpenedStream(network network.Network, stream network.Stream)        {}
func (n *Node) ClosedStream(network network.Network, stream network.Stream)        {}
