package internal

import (
	"bufio"
	"context"
	"sync"

	ggio "github.com/gogo/protobuf/io"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"

	"github.com/decanus/bureka/dht/state"
	"github.com/decanus/bureka/pb"
)

type bufferedDelimitedWriter struct {
	*bufio.Writer
	ggio.WriteCloser
}

type Writer struct {
	pool sync.Pool

	streams map[string]network.Stream

	host  host.Host
	proto protocol.ID
}

func NewWriter(h host.Host) *Writer {
	return &Writer{
		pool: sync.Pool{
			New: func() interface{} {
				w := bufio.NewWriter(nil)
				return &bufferedDelimitedWriter{
					Writer:      w,
					WriteCloser: ggio.NewDelimitedWriter(w),
				}
			},
		},
		streams: make(map[string]network.Stream),
		host:    h,
	}
}

func (w *Writer) AddStream(id state.Peer, stream network.Stream) {
	w.streams[string(id)] = stream
}

func (w *Writer) RemoveStream(id state.Peer) {
	delete(w.streams, string(id))
}

func (w *Writer) SetProtocol(proto protocol.ID) {
	w.proto = proto
}

func (w *Writer) Send(ctx context.Context, target state.Peer, msg *pb.Message) error {
	out, err := w.stream(ctx, target)
	if err != nil {
		return err
	}

	bw := w.pool.Get().(*bufferedDelimitedWriter)
	bw.Reset(out)
	err = bw.WriteMsg(msg)
	if err == nil {
		err = bw.Flush()
	}
	bw.Reset(nil)
	w.pool.Put(bw)
	return err
}

func (w *Writer) stream(ctx context.Context, target state.Peer) (network.Stream, error) {
	pid := peer.ID(target)
	out := w.streams[pid.String()]
	if out != nil {
		return out, nil
	}

	out, err := w.host.NewStream(ctx, peer.ID(target), w.proto)
	if err != nil {
		return nil, err
	}

	w.streams[pid.String()] = out
	return out, nil
}
