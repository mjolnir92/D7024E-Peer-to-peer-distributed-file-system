package kademlia

import (
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolni92/kdfs/routingtable"
)

const {
	ALPHA = 3
	K = 20
}

type T struct {
	//TODO: Create work dispatcher running <ALPHA> threads
	//TODO: Create routing table
}

func (kademlia *T) LookupContact(target *contact.T) {
	var candidates []contact.T
	//TODO: Define queried map
	closestNodes := routingtable.FindClosestContacts(target.ID, ALPHA)
	for _, node := range closestNodes {
		//TODO: Enqueue FIND_NODE call to <node> and append to <candidates>
	}
	
	//TODO: Repeat until response from k closest, remove queried from consideration
}

func (kademlia *T) LookupData(hash string) {
	// TODO
}

func (kademlia *T) Store(data []byte) {
	// TODO
}
