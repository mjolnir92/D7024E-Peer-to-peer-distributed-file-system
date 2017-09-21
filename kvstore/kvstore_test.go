package kvstore

import (
	"testing"
	"time"
	"github.com/mjolnir92/kdfs/kademliaid"
)

func TestKVStore(t *testing.T) {
	kv := New()
	data := []byte("data")
	id := kademliaid.NewHash(data)
	v := Value{timestamp: time.Now(), pin: true, data: data}

	kv.Store(v)
	//Check if the key-value pair exists in the store (we can retrieve it)
	if _, ok := kv.Get(*id); !ok {
		t.Error("TestKVStore failed, key-value pair did not get stored")
	}
}

func TestKVRemove(t *testing.T) {
	kv := New()
	data := []byte("data")
	id := kademliaid.NewHash(data)
	v := Value{timestamp: time.Now(), pin: true, data: data}

	kv.Store(v)
	kv.Remove(v)
	//Check if the key-value pair exists after removing it
	if _, ok := kv.Get(*id); ok {
		t.Error("TestKVRemove failed, key-value pair still exists")
	}
}

func TestKVGet(t *testing.T) {
	kv := New()
	data := []byte("data")
	id := kademliaid.NewHash(data)
	v := Value{timestamp: time.Now(), pin: true, data: data}

	kv.Store(v)
	//Check if the key-value pair exists in the store (we can retrieve it)
	//Can't directly compare two structs containing a []byte
	if _, ok := kv.Get(*id); !ok {
		t.Error("TestKVGet failed, couldnn't retrieve the value")
	}
}