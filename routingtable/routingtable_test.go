package routingtable

import (
	"fmt"
	"testing"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/kademliaid"
)

func TestRoutingTable(t *testing.T) {
	rt := New(contact.New(kademliaid.New("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"), 20)

	rt.AddContact(contact.New(kademliaid.New("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(contact.New(kademliaid.New("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(contact.New(kademliaid.New("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(contact.New(kademliaid.New("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(contact.New(kademliaid.New("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(contact.New(kademliaid.New("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(kademliaid.New("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}
