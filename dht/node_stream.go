package dht

import (
	"bufio"
	"context"
	"io"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-msgio"
	"github.com/pkg/errors"

	ggio "github.com/gogo/protobuf/io"
	"github.com/libp2p/go-libp2p-core/helpers"

	"github.com/decanus/bureka/pb"
)

var dhtReadMessageTimeout = 10 * time.Second
var dhtStreamIdleTimeout = 1 * time.Minute

func (n *Node) streamHandler(stream network.Stream) {
	defer stream.Reset()

	// @todo think about if we need to return bool to do stream.Close
	n.handleMessage(stream)
}

func (n *Node) handleMessage(s network.Stream) {
	r := msgio.NewVarintReaderSize(s, network.MessageSizeMax)
	ctx := n.ctx
	peer := s.Conn().RemotePeer()

	timer := time.AfterFunc(dhtStreamIdleTimeout, func() { _ = s.Reset() })
	defer timer.Stop()

	for {
		msg, err := n.latestMessage(r)
		if err != nil {
			return
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

		_, err = proto.Marshal(resp)
		if err != nil {
			// @todo
			continue
		}

		// @todo send response
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

func (n *Node) handler(t pb.Message_Type) HandlerFunc {
	switch t {
	case pb.Message_MESSAGE:
		return n.onMessage
	case pb.Message_NODE_JOIN:
		return n.onNodeJoin
	case pb.Message_NODE_ANNOUNCE:
		return n.onNodeAnnounce
	case pb.Message_NODE_EXIT:
		return n.onNodeExit
	case pb.Message_HEARTBEAT:
		return n.onHeartbeat
	case pb.Message_REPAIR_REQUEST:
		return n.onRepairRequest
	case pb.Message_STATE_REQUEST:
		return n.onStateRequest
	case pb.Message_STATE_RESPONSE:
		return n.onStateRequest
	}

	return nil
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
