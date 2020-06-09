package node

import (
	"context"
	"io"

	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-msgio"
	"github.com/pkg/errors"

	"github.com/decanus/bureka/pb"
)

func (n *Node) streamHandler(s network.Stream) {
	defer s.Reset()

	n.handleIncomingMessages(n.ctx, s)
}

func (n *Node) handleIncomingMessages(ctx context.Context, s network.Stream) {
	r := msgio.NewVarintReaderSize(s, network.MessageSizeMax)
	peer := []byte(s.Conn().RemotePeer())

	n.writer.AddStream(peer, s)

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
			logger.Infof("no handler found for message type \"%s\"", msg.Type)
			continue
		}

		resp, err := h(ctx, msg)
		if err != nil {
			logger.Error(err)
			continue
		}

		if resp == nil {
			continue
		}

		err = n.writer.Send(ctx, peer, resp)
		if err != nil {
			logger.Error(err)
		}
	}
}

func (n *Node) latestMessage(r msgio.ReadCloser) (*pb.Message, bool, error) {
	msgbytes, err := r.ReadMsg()

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
