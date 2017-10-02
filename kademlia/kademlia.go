package kademlia

import (
	"time"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/routingtable"
	"github.com/mjolnir92/kdfs/kademliaid"
)

const {
	ALPHA = 3
	K = 20
	PUBLISH_TIME = 24 * time.Hour //put these in a separate go package
	REPUBLISH_TIME = time.Hour
	EXPIRE_TIME = 24 * time.Hour
	BUCKET_REFRESH = time.Hour
}

type T struct {
	//TODO: Create work dispatcher running <ALPHA> threads, (not sure if this is the right way to do it anymore...)
	//TODO: Create routing table
	eventmanager *eventmanager.T
	kvstore *kvstore.T
	network *network.T
}

func (kademlia *T) LookupContact(target *contact.T) {
	var candidates []contact.T
	//TODO: make() candidates?
	var queried map[kademliaid.T]contact.T
	queried = make(map[kademliaid.T]contact.T)

	// Query <ALPHA> closest known nodes
	closestNodes := routingtable.FindClosestContacts(target.ID, ALPHA)
	for _, node := range closestNodes {
		//TODO: Enqueue FIND_NODE call to <node> and append to <candidates>, (or maybe simply just create a new routine?)
	}
	
	// Repeat until no closer nodes are found
	var closestSeen contact.T
	//TODO: Init <closestSeen> to closest to target in <candidates>
	newClosest := true
	for newClosest {
		// Remove queried from candidates
		//TODO: This might be wrong
		for i := 0, i < len(candidates); i++ {
			cand = candidates[i]
			if val, ok := queried[cand.ID]; ok {
				candidates = append(candidates[:i], candidates[i+1]...)
				i--
			}
		}

		//TODO: Sort <candidates> based on distance to <target> and send FIND_NODE to <ALPHA> closest? What happens if there are less than <ALPHA> unqueried candidates remaining?

	}

	//TODO: Sort <candidates> based on distance to <target> and send FIND_NODE to <K> closest. RPCs sent in batches of <ALPHA>?
	//TODO: Check if all <K> has returned?

}

func (kademlia *T) LookupData(hash string) {
	// TODO
}

func (t *T) Store(data []byte)  kademliaid.T {
	id := kademliaid.NewHash(data)
	contacts := t.LookupContact(id)
	//Defaults to the new file being unpinned
	data_val := t.kvstore.NewValue(false, data)

	for i := 0; i < len(contacts); i++ {
		go t.network.Store(contacts[i], &data_val)
	}
	//Add republish event that updates the time on the key-value pair
	f := func() {
		//If this node doesn't have the file, do LookupData to find it
		value, ok := t.kvstore.Get(id)
		if !ok {
			value = t.LookupData(id)
		}
		value.Timestamp = time.Now()

		contacts := t.LookupContact(id)
		for i := 0; i < len(contacts); i++ {
			go t.network.Store(contacts[i], &value)
		}
	}
	//Will this event ever be removed? As it looks like right now, no.
	t.eventmanager.InsertEvent(id, eventmanager.PUBLISH, f, t.PUBLISH_TIME)
}

func (kademlia *T) Cat(id Kademliaid.T) string {
	value := t.LookupData(id)
	data := value.GetData()
	return string(data[:])
}

//Updates the timestamp and sets the Pin field to true
func (kademlia *T) Pin(id Kademliaid.T) {
	//If this node doesn't have the file, do LookupData to find it
	value, ok := t.kvstore.Get(id)
	if !ok {
		value = t.LookupData(id)
	}
	value.Timestamp = time.Now()
	value.Pin = true

	contacts := t.LookupContact(id)
	for i := 0; i < len(contacts); i++ {
		go t.network.Store(contacts[i], &value)
	}
}

//Similar to Pin with the exception that the Pin field is set to false
func (kademlia *T) Unpin(id Kademliaid.T) {
	value, ok := t.kvstore.Get(id)
	if !ok {
		value = t.LookupData(id)
	}
	value.Timestamp = time.Now()
	value.Pin = false

	contacts := t.LookupContact(id)
	for i := 0; i < len(contacts); i++ {
		go t.network.Store(contacts[i], &value)
	}
}