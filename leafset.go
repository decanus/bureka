package pastry

import (
	"bytes"
	"github.com/libp2p/go-libp2p-core/peer"
	k "github.com/libp2p/go-libp2p-kbucket"
	"sort"
)

type Set []*k.PeerInfo

type LeafSet struct {
	smaller, larger Set
}

func (s Set) Closest(id peer.ID) *k.PeerInfo {
	if len(s) == 0 {
		return nil
	}

	byteid, _ := id.MarshalBinary()

	i := sort.Search(len(s), func(i int) bool {
		cmp, _ := (s)[i].Id.MarshalBinary()
		return bytes.Compare(byteid, cmp) >= 0
	})

	if i >= len(s) {
		i = len(s) - 1
	}

	return (s)[i]
}
