package kademlia

import (
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolni92/kdfs/routingtable"
)

const {
	ALPHA = 3
}

type T struct {
}

func (kademlia *T) LookupContact(target *contact.T) {
	closestNodes := routingtable.FindClosestContacts(target.ID, ALPHA)
	for _, node := range closestNodes {
		//TODO: Enqueue FIND_NODE call to <node>
	}
	//TODO: Repeat until response from k closest, remove queried from consideration
}

func (kademlia *T) LookupData(hash string) {
	// TODO
}

func (kademlia *T) Store(data []byte) {
	// TODO
}
