package bucket

import (
	"container/list"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/kademliaid"
)

type T struct {
	list *list.List
}

func New() *T {
	bucket := &T{}
	bucket.list = list.New()
	return bucket
}

func (bucket *T) AddContact(contact contact.T) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(contact.T).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
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
