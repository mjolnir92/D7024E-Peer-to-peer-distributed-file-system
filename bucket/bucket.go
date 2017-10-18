package bucket

import (
	"container/list"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/kademliaid"
)

//The front of a list is the tail of the list (most recently seen), the back of the list is the head (least recently seen)
type T struct {
	list *list.List
	replacementCache *list.List
	bucketSize int
}

func New(bucketSize int) *T {
	bucket := &T{}
	bucket.list = list.New()
	bucket.bucketSize = bucketSize
	return bucket
}

func (bucket *T) AddContact(c contact.T) {
	element := bucket.getElement(bucket.list, c)
	if element == nil {
		if bucket.list.Len() < bucket.bucketSize {
			bucket.list.PushFront(c)
		} else {
			//The bucket is full, put the contact in the replacementCache
			element = bucket.getElement(bucket.replacementCache, c)
			if element == nil {
				if bucket.replacementCache.Len() < bucket.bucketSize {
					bucket.replacementCache.PushFront(element)
				}
			} else {
				bucket.replacementCache.MoveToFront(element)
			}
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

//Remove the contact c from the bucket and replace it with the most recently seen from the replacement cache
func (bucket *T) EvictAndReplace(c contact.T) {
	element := bucket.getElement(bucket.list, c)
	if element != nil {
		if !(bucket.list.Len() < bucket.bucketSize || bucket.replacementCache.Len() == 0) {
			//If there is at least one element in the cache and the list is full, evict and replace
			bucket.list.Remove(element)
			replacement := bucket.replacementCache.Front()
			if element != nil {
				bucket.AddContact(replacement.Value.(contact.T))
				bucket.replacementCache.Remove(replacement)
			}	
		}
	}
}

func (bucket *T) getElement(l *list.List, c contact.T) *list.Element {
	var element *list.Element
	for e := l.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(contact.T).ID

		if (c).ID.Equals(nodeID) {
			element = e
		}
	}
	return element
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
