package kvstore

import (
	"github.com/mjolnir92/kdfs/kademliaid"
)

//An implementation of the Storer interface (see kvstore.go) using a map.
type Kvmap struct {
	store map[kademliaid.T]Value
}

func NewKvmap() *Kvmap {
	m := &Kvmap{}
	m.store = make(map[kademliaid.T]Value)
	return m
}

func (m *Kvmap) Get(key kademliaid.T) (Value, bool) {
	v, ok := m.store[key]
	return v, ok
}

func (m *Kvmap) Set(key kademliaid.T, v Value) {
	m.store[key] = v
}

func (m *Kvmap) Unset(key kademliaid.T) {
	delete(m.store, key)
}