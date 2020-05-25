package state_test

import (
	"testing"

	"github.com/decanus/bureka/dht/state"
)

func TestLeafSet_Insert(t *testing.T) {
	ls := state.NewLeafSet(ID())

	id := ID()

	ls.Insert(id)

	if ls.Closest(id) != id {
		t.Error("failed to insert")
	}
}

func TestLeafSet_Remove(t *testing.T) {
	ls := state.NewLeafSet(ID())

	id := ID()

	ls.Insert(id)

	if ls.Closest(id) != id {
		t.Error("failed to insert")
	}

	ls.Remove(id)

	if ls.Closest(id) != "" {
		t.Error("failed to remove")
	}
}

func TestLeafSet_Closest(t *testing.T) {
	id := ID()

	ls := state.NewLeafSet(id)

	upper := UpperID(id)
	lower := LowerID(id)

	ls.Insert(upper)
	ls.Insert(lower)

	su := UpperID(upper)
	if ls.Closest(su) != upper {
		t.Error("failed to find upper")
	}

	sl := LowerID(lower)
	if ls.Closest(sl) != lower {
		t.Error("failed to find lower")
	}
}

func TestLeafSet_Max(t *testing.T) {
	id := ID()
	u := UpperID(id)
	max := UpperID(u)

	ls := state.NewLeafSet(id)

	ls.Insert(u)
	ls.Insert(max)

	if ls.Max() != max {
		t.Error("unexpected max")
	}
}

func TestLeafSet_Min(t *testing.T) {
	id := ID()
	u := LowerID(id)
	min := LowerID(u)

	ls := state.NewLeafSet(id)

	ls.Insert(u)
	ls.Insert(min)

	if ls.Min() != min {
		t.Error("unexpected min")
	}
}

func TestLeafSet_IsInRange(t *testing.T) {
	id := ID()

	ls := state.NewLeafSet(id)
	ls.Insert(UpperID(id))
	ls.Insert(LowerID(id))

	if !ls.IsInRange(id) {
		t.Error("id not in rage as expected")
	}
}