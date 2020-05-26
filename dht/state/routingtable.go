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
	if p > len(r) {
		nr = r.grow(p)
	}

	nr[p] = nr[p].Insert(id)

	return nr
}

func (r RoutingTable) grow(n int) RoutingTable {
	nr := r
	for len(nr) <= n {
		nr = append(nr, make(Set, 0))
	}

	return nr
}

func row(self, target Peer) int {
	for i, v := range self {
		if v == target[i] {
			continue
		}

		return i - 1
	}

	return len(self) - 1
}
