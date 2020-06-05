package internal

import (
	"bufio"
	"context"
	"fmt"
	"sync"

	ggio "github.com/gogo/protobuf/io"
	"github.com/libp2p/go-libp2p-core/network"

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
}

func NewWriter() *Writer {
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
	}
}

func (w *Writer) AddStream(id state.Peer, stream network.Stream) {
	w.streams[string(id)] = stream
}

func (w *Writer) RemoveStream(id state.Peer) {
	delete(w.streams, string(id))
}

func (w *Writer) Send(ctx context.Context, target state.Peer, msg *pb.Message) error {
	out, ok := w.streams[string(target)]
	if !ok {
		return fmt.Errorf("peer %s not found", string(target))
	}

	bw := w.pool.Get().(*bufferedDelimitedWriter)
	bw.Reset(out)
	err := bw.WriteMsg(msg)
	if err == nil {
		err = bw.Flush()
	}
	bw.Reset(nil)
	w.pool.Put(bw)
	return err
}
