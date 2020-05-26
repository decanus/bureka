package state_test

import (
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/decanus/bureka/dht/state"
)

func TestSet_Insert(t *testing.T) {
	s := make(state.Set, 0)

	id := ID()

	s = s.Insert(id)

	if s.IndexOf(id) != 0 {
		t.Error("failed to insert id")
	}
}

func TestSet_Remove(t *testing.T) {
	s := make(state.Set, 0)

	id := ID()

	s = s.Insert(id)
	if s.IndexOf(id) != 0 {
		t.Error("failed to insert id")
	}

	s, ok := s.Remove(id)
	if !ok {
		t.Error("failed to remove")
	}

	if s.IndexOf(id) != -1 {
		t.Error("failed to remove id")
	}
}

func TestSet_Closest(t *testing.T) {
	s := make(state.Set, 0)

	first := ID()

	search := UpperID(first)

	second := UpperID(search)

	s = s.Insert(first)
	s = s.Insert(second)

	if !bytes.Equal(first, s.Closest(search)) {
		t.Error("unexpected closest value")
	}
}

func TestSet_Insert_IsProperlySorted(t *testing.T) {
	s := make(state.Set, 0)

	first := ID()
	second := UpperID(first)
	last := UpperID(second)

	s = s.Insert(first)
	s = s.Insert(second)
	s = s.Insert(last)

	if s.IndexOf(first) != 2 {
		t.Fatal("incorrect sorting")
	}

	if s.IndexOf(second) != 1 {
		t.Fatal("incorrect sorting")
	}

	if s.IndexOf(last) != 0 {
		t.Fatal("incorrect sorting")
	}
}

func TestSet_Insert_IsProperlySorted_Reverse(t *testing.T) {
	s := make(state.Set, 0)

	first := ID()
	second := LowerID(first)
	last := LowerID(second)

	s = s.Insert(first)
	s = s.Insert(second)
	s = s.Insert(last)

	if s.IndexOf(first) != 0 {
		t.Fatal("incorrect sorting")
	}

	if s.IndexOf(second) != 1 {
		t.Fatal("incorrect sorting")
	}

	if s.IndexOf(last) != 2 {
		t.Fatal("incorrect sorting")
	}
}

func UpperID(id state.Peer) state.Peer {
	n := make(state.Peer, len(id))
	copy(n[:], id[:])

	i := 2

	for ; i <= len(id); i++ {
		if id[i] < 255 {
			break
		}
	}

	n[i] += 1
	return n
}

func LowerID(id state.Peer) state.Peer {
	n := make(state.Peer, len(id))
	copy(n[:], id[:])

	i := 2

	for ; i <= len(id); i++ {
		if id[i] > 0 {
			break
		}
	}

	n[i] -= 1

	return n
}

func ID() []byte {
	pk, _, err := crypto.GenerateECDSAKeyPairWithCurve(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	id, err := peer.IDFromPrivateKey(pk)
	if err != nil {
		panic(err)
	}

	b, _ := id.MarshalBinary()
	return b
}
