package dht

import (
	"io"
	"time"

	proto "github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-msgio"
	"github.com/pkg/errors"

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

		go func() {
			switch msg.Type {
			case pb.Message_MESSAGE:
				n.onMessage(ctx, peer, msg)
			case pb.Message_NODE_JOIN:
				n.onNodeJoin(ctx, peer, msg)
			case pb.Message_NODE_ANNOUNCE:
				n.onNodeAnnounce(ctx, peer, msg)
			case pb.Message_NODE_EXIT:
				n.onNodeExit(ctx, peer, msg)
			case pb.Message_HEARTBEAT:
				n.onHeartbeat(ctx, peer, msg)
			case pb.Message_REPAIR_REQUEST:
				n.onRepairRequest(ctx, peer, msg)
			case pb.Message_STATE_REQUEST:
				n.onStateRequest(ctx, peer, msg)
			case pb.Message_STATE_RESPONSE:
				n.onStateRequest(ctx, peer, msg)
			}
		}()
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
