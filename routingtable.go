package routingtable

import (
	"sync"
)

const bucketSize = 20

type RoutingTable struct {
	me      Contact
	buckets [IDLength * 8]*bucket
	mux sync.Mutex
}

func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.me = me
	return routingTable
}

func (routingTable *RoutingTable) AddContact(contact Contact) {
	routingTable.mux.Lock()
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddContact(contact)
	routingTable.mux.Unlock()
}

func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) []Contact {
	routingTable.mux.Lock()
	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	defer routingTable.mux.Unlock()
	return candidates.GetContacts(count)
}

func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
	routingTable.mux.Lock()
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}
	
	defer routingTable.mux.Unlock()
	return IDLength*8 - 1
}

//Returns a pointer to the bucket closest to the target KademliaID
func (routingTable *RoutingTable) GetBucket(id *KademliaID) *bucket {
	bucketIndex := routingTable.getBucketIndex(id)
	
	routingTable.mux.Lock()
	defer routingTable.mux.Unlock()
	return routingTable.buckets[bucketIndex]
}
