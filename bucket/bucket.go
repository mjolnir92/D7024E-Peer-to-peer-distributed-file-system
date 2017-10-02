package bucket

import (
	"container/list"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/kademliaid"
)

type T struct {
	list *list.List
	bucketSize int
}

func New(bucketSize int) *T {
	bucket := &T{}
	bucket.list = list.New()
	bucket.bucketSize = bucketSize
	return bucket
}

func (bucket *T) AddContact(c contact.T) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(contact.T).ID

		if (c).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucket.bucketSize {
			bucket.list.PushFront(c)
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

func (bucket *T) GetContactAndCalcDistance(target *kademliaid.T) []contact.T {
	var contacts []contact.T

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(contact.T)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

func (bucket *T) Len() int {
	return bucket.list.Len()
}

func (bucket *T) MoveToFront(c *contact.T) {
	//bucket.list.MoveToFront(c)
}

func (bucket *T) Remove(c *contact.T) {
	
}

func (bucket *T) GetRandom() contact.T {

}