package dht

import (
	"bufio"
	"context"
	"io"

	ggio "github.com/gogo/protobuf/io"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/helpers"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-msgio"
	"github.com/pkg/errors"

	"github.com/decanus/bureka/pb"
)

func (n *Node) streamHandler(s network.Stream) {
	defer s.Reset()

	go n.handleMessageSending(n.ctx, s, n.createWriter(s.Conn().RemotePeer()))

	n.handleIncomingMessages(n.ctx, s)
}

func (n *Node) handleIncomingMessages(ctx context.Context, s network.Stream) {
	r := msgio.NewVarintReaderSize(s, network.MessageSizeMax)
	ctx := n.ctx
	peer := s.Conn().RemotePeer()

	for {
		msg, err := n.latestMessage(r)
		if err != nil {
			return
		}

		// @todo
	}

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

func (n *Node) latestMessage(r msgio.ReadCloser) (*pb.Message, error) {
	msgbytes, err := r.ReadMsg()
	// msgLen := len(msgbytes)

	if err != nil {
		r.ReleaseMsg(msgbytes)
		if err == io.EOF {
			// @todo
		}

		return nil, err
	}

	req := &pb.Message{}
	err = proto.Unmarshal(msgbytes, req)
	r.ReleaseMsg(msgbytes)
	if err != nil {
		// @todo logging?
		return nil, errors.Wrap(err, "error unmarshalling message")
	}

	return req, nil
}
