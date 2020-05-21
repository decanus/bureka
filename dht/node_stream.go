package dht

import (
	"bufio"
	"context"

	ggio "github.com/gogo/protobuf/io"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/helpers"
	"github.com/libp2p/go-libp2p-core/network"

	"github.com/decanus/bureka/pb"
)

func (n *Node) streamHandler(s network.Stream) {
	defer s.Reset()

	go n.handleMessageSending(n.ctx, s, n.createWriter(s.Conn().RemotePeer()))

	// @todo listen to messages
}

func (n *Node) handleMessageSending(ctx context.Context, s network.Stream, outgoing <-chan pb.Message) {
	bufw := bufio.NewWriter(s)
	wc := ggio.NewDelimitedWriter(bufw)

	writeMsg := func(msg proto.Message) error {
		err := wc.WriteMsg(msg)
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

			err := writeMsg(&msg)
			if err != nil {
				s.Reset()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
