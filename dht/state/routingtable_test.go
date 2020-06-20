package state_test

import (
	"bytes"
	"testing"

	"github.com/decanus/bureka/dht/state"
)

func TestRoutingTable_Insert(t *testing.T) {
	id := []byte{1, 2, 3, 4}
	insert := []byte{1, 2, 2, 3}

	r := make(state.RoutingTable, 0)

	r = r.Insert(id, insert)

	if !bytes.Equal(r[2].Get(0), insert) {
		t.Error("inserted in unexpected row")
	}
}

func TestRoutingTable_Remove(t *testing.T) {
	id := []byte{1, 2, 3, 4}
	insert := []byte{1, 2, 2, 3}

	r := make(state.RoutingTable, 0)

	r = r.Insert(id, insert)
	if !bytes.Equal(r[2].Get(0), insert) {
		t.Error("not inserted")
	}

	r = r.Remove(id, insert)
	if r[2].Length() != 0 {
		t.Error("not removed")
	}
}

func TestRoutingTable_Route(t *testing.T) {
	id := []byte{1, 2, 3, 4}
	insert := []byte{1, 2, 2, 3}
	find := []byte{1, 2, 2, 4}

	r := make(state.RoutingTable, 0)

	r = r.Insert(id, insert)

	if !bytes.Equal(r.Route(id, find), insert) {
		t.Error("unexpected route result")
	}
}

func TestRoutingTable_Route_ReturnsNoneWhenTooFar(t *testing.T) {
	id := []byte{1, 2, 3, 4}
	find := []byte{1, 2, 2, 4}

	r := make(state.RoutingTable, 0)

	if r.Route(id, find) != nil {
		t.Error("unexpected route result")
	}
}
