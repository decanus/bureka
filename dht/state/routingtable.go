package state

type RoutingTable []Set

// Route returns the node closest to the target.
func (r RoutingTable) Route(self, target Peer) Peer {
	p := commonPrefix(self, target)

	if p >= len(r) {
		// @todo error handling
		return nil
	}

	return r[p].Closest(target)
}

func (r RoutingTable) Insert(self, id Peer) RoutingTable {
	nr := r
	p := commonPrefix(self, id)
	if p > len(r) {
		nr = r.grow(p)
	}

	nr[p] = nr[p].Insert(id)

	return nr
}

func (r RoutingTable) grow(n int) RoutingTable {
	nr := r
	if n > len(r) {
		appends := len(r) - n
		for i := 0; i <= appends; i++ {
			nr = append(r, make(Set, 0))
		}
	}

	return nr
}

func commonPrefix(self, target Peer) int {
	for i, v := range self {
		if v == target[i] {
			continue
		}

		return i
	}

	return len(self)
}
