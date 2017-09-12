package contact

import (
	"fmt"
	"sort"
	"github.com/mjolnir92/kdfs/kademliaid"
)

type T struct {
	ID       *kademliaid.T
	Address  string
	distance *kademliaid.T
}

func New(id *kademliaid.T, address string) T {
	return T{id, address, nil}
}

func (contact *T) CalcDistance(target *kademliaid.T) {
	contact.distance = contact.ID.CalcDistance(target)
}

func (contact *T) Less(otherContact *T) bool {
	return contact.distance.Less(otherContact.distance)
}

func (contact *T) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}
