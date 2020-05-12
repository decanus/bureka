package state_test

import (
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/decanus/pastry/state"
)

func TestLeafSet_Insert(t *testing.T) {
	id := ID()

	ls := state.NewLeafSet(id)

	addr := Addr()

	ls.Insert(&addr)

	if ls.Closest(addr.ID).ID != addr.ID {
		t.Error("failed to insert")
	}
}

func TestLeafSet_Remove(t *testing.T) {
	id := ID()

	ls := state.NewLeafSet(id)

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

	ls := state.NewLeafSet(id)

	upper := UpperID(id)
	lower := LowerID(id)

	ls.Insert(&peer.AddrInfo{ID: upper})
	ls.Insert(&peer.AddrInfo{ID: lower})

	su := UpperID(upper)
	if ls.Closest(su).ID != upper {
		t.Error("failed to find upper")
	}

	sl := LowerID(lower)
	if ls.Closest(sl).ID != lower {
		t.Error("failed to find lower")
	}
}
