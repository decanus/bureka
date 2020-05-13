package pastry

import (
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-msgio"

	"github.com/decanus/pastry/pb"
)

var dhtReadMessageTimeout = 10 * time.Second
var dhtStreamIdleTimeout = 1 * time.Minute

func (p *Pastry) streamHandler(stream network.Stream) {
	defer stream.Reset()

	if p.handleMessage(stream) {
		stream.Close()
	}
}

func (p *Pastry) handleMessage(s network.Stream) bool {
	r := msgio.NewVarintReaderSize(s, network.MessageSizeMax)

	peer := s.Conn().RemotePeer()

	timer := time.AfterFunc(dhtStreamIdleTimeout, func() { _ = s.Reset() })
	defer timer.Stop()

	for {
		var req pb.Message
		msgbytes, err := r.ReadMsg()
		if err != nil {
			continue
		}

		msgLen := len(msgbytes)

	}

	return true
}
