package state_test

import (
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/decanus/bureka/state"
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

func TestLeafSet_Max(t *testing.T) {
	id := ID()
	u := UpperID(id)
	max := UpperID(u)

	ls := state.NewLeafSet(id)

	ls.Insert(&peer.AddrInfo{ID: u})
	ls.Insert(&peer.AddrInfo{ID: max})

	if ls.Max() != max {
		t.Error("unexpected max")
	}
}

func TestLeafSet_Min(t *testing.T) {
	id := ID()
	u := LowerID(id)
	min := LowerID(u)

	ls := state.NewLeafSet(id)

	ls.Insert(&peer.AddrInfo{ID: u})
	ls.Insert(&peer.AddrInfo{ID: min})

	if ls.Min() != min {
		t.Error("unexpected min")
	}
}

func TestLeafSet_IsInRange(t *testing.T) {
	id := ID()

	ls := state.NewLeafSet(id)
	ls.Insert(&peer.AddrInfo{ID: UpperID(id)})
	ls.Insert(&peer.AddrInfo{ID: LowerID(id)})

	if !ls.IsInRange(id) {
		t.Error("id not in rage as expected")
	}
}
