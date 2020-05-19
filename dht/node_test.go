package dht_test

import (
	"context"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	swarmt "github.com/libp2p/go-libp2p-swarm/testing"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"

	"github.com/decanus/bureka/dht"
	internal "github.com/decanus/bureka/dht/internal/mocks"
)

type MockApplication struct {
	deliver chan []byte
}

func (m *MockApplication) Deliver(msg []byte) {
	m.deliver <- msg
}

func (m *MockApplication) Forward(msg []byte, target peer.ID) bool {
	panic("implement me")
}

func (m *MockApplication) Heartbeat(id peer.ID) {
	panic("implement me")
}

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

func TestNode_Send_To_Self(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	n := setupDHT(context.Background(), t)

	mock := internal.NewMockApplication(ctrl)
	n.AddApplication(mock)

	msg := []byte("hello, world!")

	mock.EXPECT().Deliver(gomock.Eq(msg)).Times(1)

	err := n.Send(context.Background(), msg, n.ID())
	if err != nil {
		t.Fatal(err)
	}
}

func ID() peer.ID {
	pk, _, err := crypto.GenerateECDSAKeyPairWithCurve(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	id, err := peer.IDFromPrivateKey(pk)
	if err != nil {
		panic(err)
	}

	return id
}
