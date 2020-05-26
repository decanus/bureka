package state

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

func (r RoutingTable) Insert(self, id Peer) RoutingTable {
	nr := r
	p := row(self, id)
	nr = r.grow(p + 1)

	nr[p] = nr[p].Insert(id)

	return nr
}

func (r RoutingTable) grow(n int) RoutingTable {
	nr := r
	for len(nr) < n {
		nr = append(nr, make(Set, 0))
	}

	return nr
}

func (r RoutingTable) Remove(self, id Peer) RoutingTable {
	nr := r
	p := row(self, id)

	newrow, ok := nr[p].Remove(id)
	if ok {
		nr[p] = newrow
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
