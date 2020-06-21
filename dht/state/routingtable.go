package state

// RoutingTable represents a Pastry Routing table.
// Nodes are organized in rows based on the prefix of their IDs.
type RoutingTable []Set

// Route returns the node closest to the target.
func (r RoutingTable) Route(self, target Peer) Peer {
	p := row(self, target)

	if p > len(r) {
		// @todo error handling
		return nil
	}

	return r[p].Closest(target)
}

// Insert adds a Peer to the RoutingTable.
func (r RoutingTable) Insert(self, id Peer) RoutingTable {
	nr := r
	p := row(self, id)
	nr = r.grow(p + 1)

	nr[p] = nr[p].Insert(id)

	return nr
}

// Remove removes a node from the RoutingTable.
func (r RoutingTable) Remove(self, id Peer) RoutingTable {
	nr := r
	p := row(self, id)

	newrow, ok := nr[p].Remove(id)
	if ok {
		nr[p] = newrow
	}

	return nr
}

func (r RoutingTable) grow(n int) RoutingTable {
	nr := r
	for len(nr) < n {
		nr = append(nr, NewSet(10))
	}

	return nr
}

func row(self, target Peer) int {
	for i, v := range self {
		if v == target[i] {
			continue
		}

		return i
	}

	return len(self)
}
