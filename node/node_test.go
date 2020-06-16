package node

import (
	"context"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	swarmt "github.com/libp2p/go-libp2p-swarm/testing"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/decanus/bureka/dht"
	"github.com/decanus/bureka/dht/state"
	"github.com/decanus/bureka/node/internal"
)

func setupNode(ctx context.Context, t *testing.T) *Node {
	h := bhost.New(swarmt.GenSwarm(t, ctx, swarmt.OptDisableReuseport))

	w := internal.NewWriter(h)
	d := dht.New(state.Peer(h.ID()), w)

	n, err := New(ctx, d, h, w)
	if err != nil {
		t.Fatal(err)
	}

	return n
}

func setupNodes(t *testing.T, ctx context.Context, n int) []*Node {
	addrs := make([]ma.Multiaddr, n)
	dhts := make([]*Node, n)
	peers := make([]peer.ID, n)

	sanityAddrsMap := make(map[string]struct{})
	sanityPeersMap := make(map[string]struct{})

	for i := 0; i < n; i++ {
		dhts[i] = setupNode(ctx, t)
		peers[i] = peer.ID(dhts[i].dht.ID)
		addrs[i] = dhts[i].host.Addrs()[0]

		if _, lol := sanityAddrsMap[addrs[i].String()]; lol {
			t.Fatal("While setting up DHTs address got duplicated.")
		} else {
			sanityAddrsMap[addrs[i].String()] = struct{}{}
		}
		if _, lol := sanityPeersMap[peers[i].String()]; lol {
			t.Fatal("While setting up DHTs peerid got duplicated.")
		} else {
			sanityPeersMap[peers[i].String()] = struct{}{}
		}
	}

	return dhts
}

func wait(t *testing.T, ctx context.Context, a, b *Node) {
	t.Helper()

	// loop until connection notification has been received.
	// under high load, this may not happen as immediately as we would like.
	for a.dht.Find(b.dht.ID) == nil {
		select {
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		case <-time.After(time.Millisecond * 5):
		}
	}
}

func connectNoSync(t *testing.T, ctx context.Context, a, b *Node) {
	t.Helper()

	idB := peer.ID(b.dht.ID)
	addrB := b.host.Peerstore().Addrs(idB)
	if len(addrB) == 0 {
		t.Fatal("peers setup incorrectly: no local address")
	}

	a.host.Peerstore().AddAddrs(idB, addrB, (time.Minute * 2))
	pi := peer.AddrInfo{ID: idB}
	if err := a.host.Connect(ctx, pi); err != nil {
		t.Fatal(err)
	}
}

func connect(t *testing.T, ctx context.Context, a, b *Node) {
	t.Helper()
	connectNoSync(t, ctx, a, b)
	wait(t, ctx, a, b)
	wait(t, ctx, b, a)
}

func TestFindPeer(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dhts := setupNodes(t, ctx, 4)
	defer func() {
		for i := 0; i < 4; i++ {
			dhts[i].Close()
		}
	}()

	connect(t, ctx, dhts[0], dhts[1])
	connect(t, ctx, dhts[1], dhts[2])
	connect(t, ctx, dhts[1], dhts[3])

	ctxT, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	p, err := dhts[0].FindPeer(ctxT, dhts[2].host.ID())
	if err != nil {
		t.Fatal(err)
	}

	if p.ID == "" {
		t.Fatal("Failed to find peer.")
	}

	if p.ID != dhts[2].host.ID() {
		t.Fatal("Didnt find expected peer.")
	}
}
