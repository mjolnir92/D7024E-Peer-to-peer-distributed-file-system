package kademlia

import (
	"testing"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/kademliaid"
)

func TestRoutingTable(t *testing.T) {
	c0 := contact.New(kademliaid.New("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "localhost:8000")
	kademlia := New(&c0)

	c1 := contact.New(kademliaid.New("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	c2 := contact.New(kademliaid.New("00000000FFFFFFFF000000000000000000000000"), "localhost:8001")
	c3 := contact.New(kademliaid.New("0000000000000000FFFFFFFF0000000000000000"), "localhost:8001")
	c4 := contact.New(kademliaid.New("000000000000000000000000FFFFFFFF00000000"), "localhost:8001")
	c5 := contact.New(kademliaid.New("00000000000000000000000000000000FFFFFFFF"), "localhost:8001")
	c6 := contact.New(kademliaid.New("0000000000000000000000000000000000000000"), "localhost:8001")

	kademlia.routingtable.AddContact(c1)
	kademlia.routingtable.AddContact(c2)
	kademlia.routingtable.AddContact(c3)
	kademlia.routingtable.AddContact(c4)
	kademlia.routingtable.AddContact(c5)
	kademlia.routingtable.AddContact(c6)

	//FindClosestContacts should return the 5 closest contacts. c6 should not be returned
	expected := []contact.T{c1,c2,c3,c4,c5}
	contacts := kademlia.routingtable.FindClosestContacts(c0.ID, 5)
	if len(expected) != len(contacts) {
		t.Error("TestRoutingTable failed, did not return the correct amount of contacts")
	}
	for i := range contacts {
		if expected[i].ID != contacts[i].ID {
			t.Error("TestRoutingTable failed, did not return the correct contacts")
		}
	}
}