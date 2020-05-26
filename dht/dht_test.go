package dht_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/decanus/bureka/dht"
	internal "github.com/decanus/bureka/dht/internal/mocks"
	"github.com/decanus/bureka/dht/state"
	"github.com/decanus/bureka/pb"
)

func TestNode_Send_ToSelf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := internal.NewMockTransport(ctrl)
	n := dht.New([]byte("bob"), transport)

	application := internal.NewMockApplication(ctrl)
	n.AddApplication("app", application)

	msg := &pb.Message{Type: pb.Message_MESSAGE, Key: n.ID}

	application.EXPECT().Deliver(gomock.Eq(msg)).Times(1)

	err := n.Send(context.Background(), msg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNode_Send_WhenPeerInLeafSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := internal.NewMockTransport(ctrl)
	n := dht.New([]byte("bob"), transport)

	application := internal.NewMockApplication(ctrl)
	n.AddApplication("app", application)

	target := make(state.Peer, 3)
	target[0] = 3
	n.AddPeer(target)

	msg := &pb.Message{Type: pb.Message_MESSAGE, Key: target}

	application.EXPECT().Forward(gomock.Eq(msg), gomock.Eq(target)).Times(1).Return(true)

	ctx := context.Background()
	transport.EXPECT().Send(gomock.Eq(ctx), gomock.Eq(target), gomock.Eq(msg)).Times(1)

	err := n.Send(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNode_Send_DoesNothingOnFalseForward(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := internal.NewMockTransport(ctrl)
	n := dht.New([]byte("bob"), transport)

	application := internal.NewMockApplication(ctrl)
	n.AddApplication("app", application)

	target := make(state.Peer, 3)
	target[0] = 3
	n.AddPeer(target)

	msg := &pb.Message{Type: pb.Message_MESSAGE, Key: target}

	application.EXPECT().Forward(gomock.Eq(msg), gomock.Eq(target)).Times(1).Return(false)

	ctx := context.Background()

	err := n.Send(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
}
