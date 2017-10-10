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

//The ping callback function should be the ping method from kademlia/network
func (bucket *T) AddContact(c contact.T, ping func(c2 *contact.T) error) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(contact.T).ID

		if (c).ID.Equals(nodeID) {
			element = e
		}
	}

	//The front of the list is the tail of the list, the back of the list is the head
	if element == nil {
		if bucket.list.Len() < bucket.bucketSize {
			bucket.list.PushFront(c)
		} else {
			//ping least-recently seen node, evict if unresponsive and insert the new contact at the tail
			//if it responds move it to the front of the list and discard the new contact
			leastRecent := bucket.list.Back()
			if leastRecent != nil {
				leastRecentContact := leastRecent.Value.(contact.T)
				if err := ping(&leastRecentContact); err != nil {
					bucket.list.Remove(leastRecent)
					bucket.list.PushFront(c)
				} else {
					bucket.list.MoveToFront(leastRecent)
				}
			}
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
