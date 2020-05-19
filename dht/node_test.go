package dht_test

import (
	"context"
	"crypto/elliptic"
	"crypto/rand"
	"reflect"
	"testing"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	swarmt "github.com/libp2p/go-libp2p-swarm/testing"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"

	"github.com/decanus/bureka/dht"
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
	n := setupDHT(context.Background(), t)

	m := &MockApplication{deliver: make(chan []byte)}
	n.AddApplication(m)

	msg := []byte("hello, world!")

	go func() {
		err := n.Send(context.Background(), msg, n.ID())
		if err != nil {
			t.Fatal(err)
		}
	}()

	b := <- m.deliver

	if !reflect.DeepEqual(b, msg) {
		t.Errorf("expected: %v actual: %v", msg, b)
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