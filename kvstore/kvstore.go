package kvstore

import (
	"sync"
	"time"
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

//Function to store a key-value pair
func (t *T) Store(v Value) {
	t.mux.Lock()
	//Create a kademliaid (key) for the value to be inserted.
	data := v.GetData()
	key := kademliaid.NewHash(data)

	current, ok := t.store.Get(*key)
	if ok {
		//The key did exist
		if current.Before(v) {
			t.store.Set(*key, v)
		}
	} else {
		//Key did not already exist
		t.store.Set(*key, v)
	}
	t.mux.Unlock()
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

//Pin update
//Updates pin if the timestamp is newer
func (t *T) Pin(id kademliaid.T, time time.Time) {
	//TODO
}

func (t *T) Get(key kademliaid.T) (Value, bool) {
	v, ok := t.store.Get(key)
	return v, ok
}