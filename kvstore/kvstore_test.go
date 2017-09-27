package kvstore

import (
	"bytes"
	"testing"
	"github.com/mjolnir92/kdfs/kademliaid"
)

func TestKVStore(t *testing.T) {
	kv := New()
	data := []byte("data")
	id := kademliaid.NewHash(data)
	v := NewValue(true, data)

	kv.Store(v)
	got, ok := kv.Get(*id)
	//Check if the key-value pair exists in the store (we can retrieve it)
	if !ok {
		t.Error("TestKVStore failed, key-value pair did not get stored")
	}
	if !bytes.Equal(got.GetData(), data) {
		t.Error("TestKVStore failed, wrong value")
	}
	kv.Remove(v)
	//Check if the key-value pair exists after removing it
	if _, ok := kv.Get(*id); ok {
		t.Error("TestKVStore failed, key-value pair did not get removed")
	}
}