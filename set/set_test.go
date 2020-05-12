package set_test

import (
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/decanus/pastry/set"
)

func TestSet_Insert(t *testing.T) {
	s := make(set.Set, 0)

	addr := Addr()

	s = s.Insert(&addr)

	if s.IndexOf(addr.ID) != 0 {
		t.Error("failed to insert id")
	}
}

func TestSet_Remove(t *testing.T) {
	s := make(set.Set, 0)

	addr := Addr()

	s = s.Insert(&addr)
	if s.IndexOf(addr.ID) != 0 {
		t.Error("failed to insert id")
	}

	s, ok := s.Remove(addr.ID)
	if !ok {
		t.Error("failed to remove")
	}

	if s.IndexOf(addr.ID) != -1 {
		t.Error("failed to remove id")
	}
}

func TestSet_Closest(t *testing.T) {
	s := make(set.Set, 0)

	first := Addr()

	bytes, _ := first.ID.MarshalBinary()
	bytes[2] += 1

	search, _ := peer.IDFromBytes(bytes)

	bytes[2] += 1
	second, _ := peer.IDFromBytes(bytes)

	s = s.Insert(&first)
	s = s.Insert(&peer.AddrInfo{ID: second})

	if first.ID != s.Closest(search).ID {
		t.Error("unexpected closest value")
	}
}

func Addr() peer.AddrInfo {
	return peer.AddrInfo{
		ID: ID(),
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
