package kvstore

import (
	"sync"
	"github.com/mjolnir92/kdfs/kademliaid"
)

type T struct{
	store storer
	mux sync.Mutex
}

type storer interface{
	Get(kademliaid.T) (Value, bool)
	Set(kademliaid.T, Value)
	Unset(kademliaid.T)
}

//Defaults to creating a T with a kvmap.
//TODO: Options for selecting Store implementation
func New() *T {
	t := &T{}
	t.store = NewKvmap()
	return t
}

//Function to store a key-value pair. Returns true if the value was inserted
func (t *T) Store(v Value) bool {
	t.mux.Lock()
	//Create a kademliaid (key) for the value to be inserted.
	data := v.GetData()
	key := kademliaid.NewHash(data)
	inserted := false

	current, ok := t.store.Get(*key)
	if ok {
		//The key did exist
		if current.Before(v) {
			t.store.Set(*key, v)
			inserted = true
		}
	} else {
		//Key did not already exist
		t.store.Set(*key, v)
		inserted = true
	}
	t.mux.Unlock()
	return inserted
}

//Removes a key-value pair from the storer
func (t *T) Remove(v Value) {
	t.mux.Lock()
	//Create a kademliaid (key) for the value to be inserted.
	data := v.GetData()
	key := kademliaid.NewHash(data)

	_, ok := t.store.Get(*key)
	if ok {
		t.store.Unset(*key)
	}
	t.mux.Unlock()
}

func (t *T) Get(key kademliaid.T) (Value, bool) {
	t.mux.Lock()
	v, ok := t.store.Get(key)
	t.mux.Unlock()
	return v, ok
}