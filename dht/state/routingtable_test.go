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

	if !bytes.Equal(r[1][0], insert) {
		t.Error("inserted in unexpected row")
	}
}
