package set_test

import (
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/decanus/pastry/set"
)

func TestLeafSet_Insert(t *testing.T) {
	id := ID()

	ls := set.NewLeafSet(id)

	addr := Addr()

	ls.Insert(&addr)

	if ls.Closest(addr.ID).ID != addr.ID {
		t.Error("failed to insert")
	}
}

func TestLeafSet_Remove(t *testing.T) {
	id := ID()

	ls := set.NewLeafSet(id)

	addr := Addr()

	ls.Insert(&addr)

	if ls.Closest(addr.ID).ID != addr.ID {
		t.Error("failed to insert")
	}

	ls.Remove(addr.ID)

	if ls.Closest(addr.ID) != nil {
		t.Error("failed to remove")
	}
}

func TestLeafSet_Closest(t *testing.T) {
	id := ID()

	ls := set.NewLeafSet(id)

	b, _ := id.MarshalBinary()

	ub := make([]byte, len(b))
	copy(ub, b)
	ub[2] += 1

	lb := make([]byte, len(b))
	copy(lb, b)
	lb[2] -= 1

	upper, _ := peer.IDFromBytes(ub)
	lower, _ := peer.IDFromBytes(lb)

	ls.Insert(&peer.AddrInfo{ID: upper})
	ls.Insert(&peer.AddrInfo{ID: lower})

	ub[2] += 1
	su, _ := peer.IDFromBytes(ub)

	if ls.Closest(su).ID != upper {
		t.Error("failed to find upper")
	}

	lb[2] -= 1
	sl, _ := peer.IDFromBytes(lb)

	if ls.Closest(sl).ID != lower {
		t.Error("failed to find lower")
	}
}
