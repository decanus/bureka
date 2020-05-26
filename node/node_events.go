package node

import (
	"github.com/libp2p/go-eventbus"
	"github.com/libp2p/go-libp2p-core/event"
)

func (n *Node) subscribe() (event.Subscription, error) {
	bufSize := eventbus.BufSize(256)

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

	s, err := n.host.EventBus().Subscribe(evts, bufSize)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (n *Node) poll(s event.Subscription) {
	defer s.Close()
}
