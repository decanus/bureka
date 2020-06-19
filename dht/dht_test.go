package dht_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/decanus/bureka/dht"
	internal "github.com/decanus/bureka/dht/internal/mocks"
	"github.com/decanus/bureka/dht/state"
	"github.com/decanus/bureka/pb"
)

func TestDHT_AddPeer_And_RemovePeer(t *testing.T) {
	id := []byte{5, 5, 5, 5}
	insert := []byte{0, 1, 3, 3}
	d := dht.New(id)

	d.AddPeer(insert)

	if !bytes.Equal(d.LeafSet.Closest(insert), insert) {
		t.Error("failed to insert in LeafSet")
	}

	if !bytes.Equal(d.NeighborhoodSet.Closest(insert), insert) {
		t.Error("failed to insert in NeighborhoodSet")
	}

	if !bytes.Equal(d.RoutingTable.Route(id, insert), insert) {
		t.Error("failed to insert in NeighborhoodSet")
	}

	d.RemovePeer(insert)

	if d.RoutingTable.Route(id, insert) != nil {
		t.Error("failed to remove peer from RoutingTable")
	}

	if d.NeighborhoodSet.Closest(insert) != nil {
		t.Error("failed to remove peer from NeighborhoodSet")
	}

	if d.LeafSet.Closest(insert) != nil {
		t.Error("failed to remove peer from LeafSet")
	}
}

func TestDHT_Send_ToSelf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := dht.New([]byte("bob"))

	application := internal.NewMockApplication(ctrl)
	d.AddApplication("app", application)

	msg := &pb.Message{Type: pb.Message_MESSAGE, Key: d.ID}

	application.EXPECT().Deliver(gomock.Eq(msg)).Times(1)

	err := d.Send(context.Background(), msg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDHT_Send_WhenPeerInLeafSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := dht.New([]byte("bob"))

	application := internal.NewMockApplication(ctrl)
	d.AddApplication("app", application)

	target := make(state.Peer, 3)
	target[0] = 3
	d.AddPeer(target)

	msg := &pb.Message{Type: pb.Message_MESSAGE, Key: target}

	application.EXPECT().Forward(gomock.Eq(msg), gomock.Eq(target)).Times(1).Return(true)

	ctx := context.Background()

	c := make(chan dht.Packet, 5)
	d.Feed().Subscribe(c)

	err := d.Send(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}

	val := <-c
	if msg != val.Message {
		t.Fail()
	}
}

func TestDHT_Send_DoesNothingOnFalseForward(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := dht.New([]byte("bob"))

	application := internal.NewMockApplication(ctrl)
	d.AddApplication("app", application)

	target := make(state.Peer, 3)
	target[0] = 3
	d.AddPeer(target)

	msg := &pb.Message{Type: pb.Message_MESSAGE, Key: target}

	application.EXPECT().Forward(gomock.Eq(msg), gomock.Eq(target)).Times(1).Return(false)

	ctx := context.Background()

	err := d.Send(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
}
