package d7024e

import (
	"config"
	"contact"
	"routingtable"
)

const {
	ALPHA = 3
}

type Kademlia struct {
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	closestNodes := routingtable.FindClosestContacts(target.ID, ALPHA)
	for _, node := range closestNodes {
		//TODO: Enqueue FIND_NODE call to <node>
	}
	//TODO: Repeat until response from k closest, remove queried from consideration
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
