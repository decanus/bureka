package dht

import (
	"bufio"
	"context"

	"github.com/libp2p/go-libp2p-core/helpers"
	"github.com/libp2p/go-libp2p-core/network"
)

func (n *Node) handleMessageSending(ctx context.Context, s network.Stream, outgoing <-chan []byte) {
	bufw := bufio.NewWriter(s)

	writeMsg := func(msg []byte) error {
		_, err := bufw.Write(msg)
		if err != nil {
			return err
		}

		return bufw.Flush()
	}

	defer helpers.FullClose(s)
	for {
		select {
		case msg, ok := <-outgoing:
			if !ok {
				return
			}

			err := writeMsg(msg)
			if err != nil {
				s.Reset()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

