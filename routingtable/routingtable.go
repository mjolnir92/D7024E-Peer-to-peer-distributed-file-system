package routingtable

import (
	"sync"
	"sort"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/bucket"
	"github.com/mjolnir92/kdfs/kademliaid"
)

const bucketSize = 20

type T struct {
	me      contact.T
	buckets [kademliaid.IDLength * 8]*bucket.T
	mux sync.Mutex
}

func New(me contact.T, bucketSize int) *T {
	routingTable := &T{}
	for i := 0; i < kademliaid.IDLength*8; i++ {
		routingTable.buckets[i] = bucket.New(bucketSize)
	}
	routingTable.me = me
	return routingTable
}

func (routingTable *T) AddContact(contact contact.T) {
	routingTable.mux.Lock()
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddContact(contact)
	routingTable.mux.Unlock()
}

func (routingTable *T) FindClosestContacts(target *kademliaid.T, count int) []contact.T {
	routingTable.mux.Lock()
	var candidates []contact.T
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates = append(candidates, bucket.GetContactAndCalcDistance(target)...)

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < kademliaid.IDLength*8) && len(candidates) < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates = append(candidates, bucket.GetContactAndCalcDistance(target)...)
		}
		if bucketIndex+i < kademliaid.IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates = append(candidates, bucket.GetContactAndCalcDistance(target)...)
		}
	}

	sort.Sort(contact.ByDist(candidates))

	if count > len(candidates) {
		count = len(candidates)
	}

	defer routingTable.mux.Unlock()
	return candidates[:count]
}

func (routingTable *T) FindKClosestContacts(target *kademliaid.T) []contact.T {
	return routingTable.FindClosestContacts(target, bucketSize)
}

func (routingTable *T) getBucketIndex(id *kademliaid.T) int {
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < kademliaid.IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}
	
	return kademliaid.IDLength*8 - 1
}

//Returns a pointer to the bucket closest to the target KademliaID
func (routingTable *T) GetBucket(id *kademliaid.T) *bucket.T {
	bucketIndex := routingTable.getBucketIndex(id)
	
	routingTable.mux.Lock()
	defer routingTable.mux.Unlock()
	return routingTable.buckets[bucketIndex]
}
