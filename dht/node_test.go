package dht_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	swarmt "github.com/libp2p/go-libp2p-swarm/testing"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"

	"github.com/decanus/bureka/dht"
	internal "github.com/decanus/bureka/dht/internal/mocks"
	"github.com/decanus/bureka/pb"
)

// @TODO MORE TESTS

func setupDHT(ctx context.Context, t *testing.T) *dht.Node {
	d, err := dht.New(
		ctx,
		bhost.New(swarmt.GenSwarm(t, ctx, swarmt.OptDisableReuseport)),
	)

	if err != nil {
		t.Fatal(err)
	}

	return d
}

func connectNoSync(t *testing.T, ctx context.Context, a, b *dht.Node) {
	t.Helper()

	idB := b.ID()
	addrB := b.Host.Peerstore().Addrs(idB)
	if len(addrB) == 0 {
		t.Fatal("peers setup incorrectly: no local address")
	}

	a.Host.Peerstore().AddAddrs(idB, addrB, peerstore.TempAddrTTL)
	pi := peer.AddrInfo{ID: idB}
	if err := a.Host.Connect(ctx, pi); err != nil {
		t.Fatal(err)
	}
}

func TestNode_Send_To_Self(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	n := setupDHT(context.Background(), t)

	mock := internal.NewMockApplication(ctrl)
	n.AddApplication(mock)

	msg := pb.Message{ Type: pb.Message_MESSAGE, Key: string(n.ID())}

	mock.EXPECT().Deliver(gomock.Eq(msg)).Times(1)

	err := n.Send(context.Background(), msg)
	if err != nil {
		t.Fatal(err)
	}
}

