package contact

import (
	"fmt"
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

// This implements sort.Interface for []T based on the distance
type ByDist []T

func (a ByDist) Len() int {
	return len(a)
}

func (a ByDist) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByDist) Less(i, j int) bool {
	return a[i].Less(&a[j])
}
