package dht_test

import (
	"bytes"
	"testing"

	"github.com/decanus/bureka/dht"
)

func TestNode_AddPeer_And_RemovePeer(t *testing.T) {
	id := []byte{5, 5, 5, 5}
	insert := []byte{0, 1, 3, 3}
	n := dht.New(id)

	n.AddPeer(insert)

	if !bytes.Equal(n.LeafSet.Closest(insert), insert) {
		t.Error("failed to insert in LeafSet")
	}

	if !bytes.Equal(n.NeighborhoodSet.Closest(insert), insert) {
		t.Error("failed to insert in NeighborhoodSet")
	}

	if !bytes.Equal(n.RoutingTable.Route(id, insert), insert) {
		t.Error("failed to insert in NeighborhoodSet")
	}

	n.RemovePeer(insert)

	if n.RoutingTable.Route(id, insert) != nil {
		t.Error("failed to remove peer from RoutingTable")
	}

	if n.NeighborhoodSet.Closest(insert) != nil {
		t.Error("failed to remove peer from NeighborhoodSet")
	}

	if n.LeafSet.Closest(insert) != nil {
		t.Error("failed to remove peer from LeafSet")
	}
}

