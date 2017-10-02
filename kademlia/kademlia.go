package kademlia

import (
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolni92/kdfs/routingtable"
	"github.com/mjolni92/kdfs/kademliaid"
)

const {
	ALPHA = 3
	K = 20
}

type T struct {
	//TODO: Create work dispatcher running <ALPHA> threads, (not sure if this is the right way to do it anymore...)
	//TODO: Create routing table
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

func (kademlia *T) Store(data []byte) {
	// TODO
}
