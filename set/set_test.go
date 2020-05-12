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

	search := UpperID(first.ID)

	second := UpperID(search)

	s = s.Insert(&first)
	s = s.Insert(&peer.AddrInfo{ID: second})

	if first.ID != s.Closest(search).ID {
		t.Error("unexpected closest value")
	}
}

func TestSet_Insert_IsProperlySorted(t *testing.T) {
	s := make(set.Set, 0)

	first := ID()
	second := UpperID(first)
	last := UpperID(second)

	s = s.Insert(&peer.AddrInfo{ID: first})
	s = s.Insert(&peer.AddrInfo{ID: second})
	s = s.Insert(&peer.AddrInfo{ID: last})

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
	s := make(set.Set, 0)

	first := ID()
	second := LowerID(first)
	last := LowerID(second)

	s = s.Insert(&peer.AddrInfo{ID: first})
	s = s.Insert(&peer.AddrInfo{ID: second})
	s = s.Insert(&peer.AddrInfo{ID: last})

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

func UpperID(id peer.ID) peer.ID {
	b, _ := id.MarshalBinary()

	i := 2

	for ; i <= len(b); i++ {
		if b[i] < 255 {
			break
		}
	}

	b[i] += 1

	p, _ := peer.IDFromBytes(b)

	return p
}

func LowerID(id peer.ID) peer.ID {
	b, _ := id.MarshalBinary()
	i := 2

	for ; i <= len(b); i++ {
		if b[i] > 0 {
			break
		}
	}

	b[i] -= 1

	p, _ := peer.IDFromBytes(b)

	return p
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
