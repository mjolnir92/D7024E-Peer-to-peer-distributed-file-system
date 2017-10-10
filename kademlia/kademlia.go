package kademlia

import (
	"net"
	"time"
	"sort"
	"sync"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/routingtable"
	"github.com/mjolnir92/kdfs/kademliaid"
	"github.com/mjolnir92/kdfs/constants"
	"github.com/mjolnir92/kdfs/eventmanager"
	"github.com/mjolnir92/kdfs/kvstore"
)

type Candidates struct {
	c	[]contact.T
	mux	sync.Mutex
}

func (candidates *Candidates) CalcDistances(target *kademliaid.T) {
	for _, c := range candidates.c {
		c.CalcDistance(target)
	}
}

type T struct {
	eventmanager *eventmanager.T
	kvstore *kvstore.T
	routingtable *routingtable.T
	contactMe *contact.T
	conn *net.UDPConn
}

func New(contactMe *contact.T) *T{
	t := &T{}
	t.contactMe = contactMe
	t.eventmanager = eventmanager.New()
	t.routingtable = routingtable.New(*t.contactMe, t.eventmanager, constants.K, t.PingNoRefresh)
	t.kvstore = kvstore.New()

	for i := 0; i < kademliaid.IDLength*8; i++{
		f := func() {
			t.refreshBucket(i)
		}
		t.eventmanager.InsertEvent(*contactMe.ID, i, f, constants.BUCKET_REFRESH)
	}
	return t
}

//This method refreshes the bucket corresponding to the index
func (t *T) refreshBucket(index int) {
	randomID := kademliaid.NewRandomCommonPrefix(*t.contactMe.ID, uint8(index))
	contacts := t.LookupContact(randomID)
	for _, c := range(contacts) {
		t.routingtable.AddContact(c)
	}
}

func (t *T) Join(c *contact.T) {
	t.routingtable.AddContact(*c)
	//Does LookupContact send rpcs to all returned contacts? If so i dont need to add them to the routingtable heres
	contacts := t.LookupContact(t.contactMe.ID)
	for _, c := range(contacts) {
		t.routingtable.AddContact(c)
	}
	//Refresh all buckets further away than the closest neighbor
	closestNeighbors := t.routingtable.FindClosestContacts(t.contactMe.ID, 1)
	closestNeighbor := closestNeighbors[0]
	index := t.routingtable.GetBucketIndex(closestNeighbor.ID)
	for i := index; i < kademliaid.IDLength; i++ {
		t.refreshBucket(i)
		t.eventmanager.ResetEvent(*t.contactMe.ID, i, constants.BUCKET_REFRESH)
	}
}

func (t *T) issueFindNode(node *contact.T, target *kademliaid.T, candidates *Candidates, i int,  queried map[kademliaid.T]contact.T, replied map[kademliaid.T]contact.T, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := t.FindNode(node, target)
	queried[*node.ID] = *node
	candidates.mux.Lock()
	if err != nil {
		if i != -1 {
			candidates.c = append(candidates.c[:i], candidates.c[i+1:]...)
		}
	} else {
		replied[*node.ID] = *node
		candidates.c = append(candidates.c, res...)
		candidates.CalcDistances(target)
		sort.Sort(contact.ByDist(candidates.c))
	}
	candidates.mux.Unlock()
}

func (t *T) LookupContact(target *kademliaid.T) []contact.T {
	candidates := Candidates{c: make([]contact.T, 0)}
	queried := make(map[kademliaid.T]contact.T)
	replied := make(map[kademliaid.T]contact.T)
	var wg sync.WaitGroup
	/*
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
	}()
	*/

	// Query <ALPHA> closest known nodes
	closestNodes := t.routingtable.FindClosestContacts(target, constants.ALPHA)
	for _, node := range closestNodes {
		/*
		go func() {
			res, err := t.FindNode(node, target)
			queried[node.ID] = node
			if err != nil {
				//TODO: Handle error
			} else {
				replied[node.ID] = node
				ch <- res
			}
		}()
		*/
		// Call with i = -1 do denote that there is nothing to evict from candidates yet
		go t.issueFindNode(&node, target, &candidates, -1, queried, replied, &wg)
	}

	// Repeat until no closer nodes are found
	for {
		wg.Add(constants.ALPHA)
		aCount := 0
		candidates.mux.Lock()
		closestSeen := candidates.c[0]
		for i := 0; i < constants.K; i++ {
			if _, ok := queried[*candidates.c[i].ID]; !ok {
				/*
				go func() {
					res, err := t.FindNode(candidates.c[i], target)
					queried[candidates.c[i].ID] = candidates.c[i]
					if err != nil {
						//TODO: Handle error
					} else {
						replied[candidates.c[i].ID] = candidates.c[i]
						ch <- res
					}
				}()
				*/
				go t.issueFindNode(&candidates.c[i], target, &candidates, i, queried, replied, &wg)
				aCount++
			}
			if aCount >= constants.ALPHA {
				break
			}
		}
		candidates.mux.Unlock()

		//  Wait for responses
		wg.Wait()

		candidates.mux.Lock()
		if closestSeen.ID == candidates.c[0].ID {
			break
		}
		candidates.mux.Unlock()
	}

	pendingReplies := true
	// Query all K closest candidates that have not been queried until all have responded
	for pendingReplies {
		pendingReplies = false
		candidates.mux.Lock()
		for i := 0; i < constants.K; i++ {
			if _, ok := queried[*candidates.c[i].ID]; !ok {
				/*
				go func() {
					res, err := t.FindNode(candidates.c[i], target)
					queried[candidates.c[i].ID] = candidates.c[i]
					if err != nil {
						//TODO: Handle error
					} else {
						replied[candidates.c[i].ID] = candidates.c[i]
						ch <- res
					}
				}( i,)
				*/
				go t.issueFindNode(&candidates.c[i], target, &candidates, i, queried, replied, &wg)
			}
		}
		candidates.mux.Unlock()

		candidates.mux.Lock()
		for i := 0; i < constants.K; i++ {
			if _, ok := replied[*candidates.c[i].ID]; !ok {
				pendingReplies = true
			}
		}
		candidates.mux.Unlock()
	}

	candidates.mux.Lock()
	defer candidates.mux.Unlock()
	return candidates.c[:constants.K]
}

func (t *T) LookupData(target *kademliaid.T) kvstore.Value {
	// TODO
	//remove whats below this comment, it was just to remove compile errors temporarily
	var data kvstore.Value
	return data
}

func (t *T) KademliaStore(data []byte)  kademliaid.T {
	id := kademliaid.NewHash(data)
	contacts := t.LookupContact(id)
	//Defaults to the new file being unpinned
	data_val := kvstore.NewValue(false, data)

	for i := 0; i < len(contacts); i++ {
		go t.Store(&contacts[i], &data_val)
	}
	//Add republish event that updates the time on the key-value pair
	f := func() {
		//If this node doesn't have the file, do LookupData to find it
		value, ok := t.kvstore.Get(*id)
		if !ok {
			value = t.LookupData(id)
		}
		value.Timestamp = time.Now()

		contacts := t.LookupContact(id)
		for i := 0; i < len(contacts); i++ {
			go t.Store(&contacts[i], &value)
		}
	}
	//Will this event ever be removed? As it looks like right now, no.
	t.eventmanager.InsertEvent(*id, constants.PUBLISH, f, constants.PUBLISH_TIME)
	return *id
}

func (t *T) Cat(id kademliaid.T) string {
	value := t.LookupData(&id)
	return string(value.GetData()[:])
}

//Updates the timestamp and sets the Pin field to true
func (t *T) Pin(id kademliaid.T) {
	//If this node doesn't have the file, do LookupData to find it
	value, ok := t.kvstore.Get(id)
	if !ok {
		value = t.LookupData(&id)
	}
	value.Timestamp = time.Now()
	value.Pin = true

	contacts := t.LookupContact(&id)
	for i := 0; i < len(contacts); i++ {
		go t.Store(&contacts[i], &value)
	}
}

//Similar to Pin with the exception that the Pin field is set to false
func (t *T) Unpin(id kademliaid.T) {
	value, ok := t.kvstore.Get(id)
	if !ok {
		value = t.LookupData(&id)
	}
	value.Timestamp = time.Now()
	value.Pin = false

	contacts := t.LookupContact(&id)
	for i := 0; i < len(contacts); i++ {
		go t.Store(&contacts[i], &value)
	}
}
