package kademlia

import (
	"time"
	"sort"
	"sync"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/routingtable"
	"github.com/mjolnir92/kdfs/kademliaid"
	"github.com/mjolnir92/kdfs/constants"
	"github.com/mjolnir92/kdfs/network"
)

const (
	ALPHA = 3
	K = 20
)

type Candidates struct {
	c	[]contact.T
	mux	sync.Mutex
}

func (candidates *Candidates) CalcDistances(target *contact.T) {
	for _, c := range candidates.c {
		c.CalcDistance(target.ID)
	}
}

type T struct {
	eventmanager *eventmanager.T
	kvstore *kvstore.T
	network *network.T
}

func (t *T) LookupContact(target *contact.T) []contact.T {
	candidates := Candidates{c: make([]contact.T, 0)}
	queried := make(map[kademliaid.T]contact.T)
	ch := make(chan []contact.T)

	// Routine that updates candidates
	go func() {
		for i := range ch {
			candidates.mux.Lock()
			candidates.c = append(candidates.c, i...)
			candidates.CalcDistances(target)
			sort.Sort(contact.ByDist(candidates.c))
			candidates.mux.Unlock()
		}
	}

	// Query <ALPHA> closest known nodes
	closestNodes := t.network.routingtable.FindClosestContacts(target.ID, ALPHA)
	for _, node := range closestNodes {
		go func() {
			res, err := t.network.FindNode(node, target.ID)
			if err != nil {
				//TODO: Handle error
			} else {
				queried[node.ID] = node
				ch <- res
			}
		}
	}

	// Repeat until no closer nodes are found
	for {
		closestSeen := candidates.c[0]
		aCount := 0
		for i := 0; i < K; i++ {
			if val, ok := queried[candidates.c[i]]; !ok {
				go func() {
					res, err := t.network.FindNode(node, target.ID)
					if err != nil {
						//TODO: Handle error
					} else {
						queried[node.ID] = node
						ch <- res
					}
				}
				aCount++
			}
			if aCount >= ALPHA {
				break
			}
		}

		if closestSeen.ID == candidates.c[0].ID {
			break
		}
	}

	pendingReplies := true
	// Query all K closest candidates that have not been queried until all have responded
	for pendingReplies {
		pendingReplies = false
		for i := 0; i < K; i++ {
			if val, ok := queried[candidates.c[i]]; !ok {
				go func() {
					res, err := t.network.FindNode(node, target.ID)
					if err != nil {
						//TODO: Handle error
					} else {
						queried[node.ID] = node
						ch <- res
					}
				}
				pendingReplies = true
			}
		}
	}

	return candidates.c[:K]
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
	t.eventmanager.InsertEvent(id, eventmanager.PUBLISH, f, constants.PUBLISH_TIME)
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
