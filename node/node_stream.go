package node

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

	//go n.handleMessageSending(n.ctx, s, n.createWriter(s.Conn().RemotePeer()))

	n.handleIncomingMessages(n.ctx, s)
}

func (n *Node) handleIncomingMessages(ctx context.Context, s network.Stream) {
	r := msgio.NewVarintReaderSize(s, network.MessageSizeMax)
	peer, _ := s.Conn().RemotePeer().MarshalBinary()

	for {
		msg, done, err := n.latestMessage(r)
		if done {
			return
		}

		if err != nil {
			continue // @todo?
		}

		h := n.handler(msg.Type)
		if h == nil {
			// @todo
			continue
		}

		resp := h(ctx, peer, msg)
		if resp == nil {
			// @todo
			continue
		}

		// @todo send to our output
		//err = n.Send(ctx, peer, *resp)
		//if err != nil {
		//	// @todo
		//}
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

func (n *Node) latestMessage(r msgio.ReadCloser) (*pb.Message, bool, error) {
	msgbytes, err := r.ReadMsg()
	// msgLen := len(msgbytes)

	if err != nil {
		r.ReleaseMsg(msgbytes)
		if err == io.EOF {
			return nil, true, nil
			// @todo
		}

		return nil, false, err
	}

	req := &pb.Message{}
	err = proto.Unmarshal(msgbytes, req)
	r.ReleaseMsg(msgbytes)
	if err != nil {
		// @todo logging?
		return nil, false, errors.Wrap(err, "error unmarshalling message")
	}

	return req, false, nil
}
