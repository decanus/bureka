package pastry

import "github.com/libp2p/go-libp2p-core/network"

func (p *Pastry) streamHandler(stream network.Stream) {
	defer stream.Reset()

	if p.handleMessage(stream) {
		stream.Close()
	}
}

func (p *Pastry) handleMessage(s network.Stream) bool {
	return true
}
