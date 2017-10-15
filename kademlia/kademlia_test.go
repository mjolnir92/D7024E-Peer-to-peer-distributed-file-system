package kademlia

import (
	"log"
	"time"
	"testing"
	"github.com/mjolnir92/kdfs/kademliaid"
	"github.com/mjolnir92/kdfs/contact"
)

func TestLookupData(t *testing.T) {
	address := "localhost:12300"
	ct_client := contact.New(kademliaid.New("1000000000000000000000000000000000000000"), address)
	nw_client := New(&ct_client)
	go nw_client.Listen("localhost", 12300)
	time.Sleep(50 * time.Millisecond)
	
	ct_server1 := contact.New(kademliaid.New("1100000000000000000000000000000000000000"), "localhost:12301")
	nw_server1 := New(&ct_server1)
	go nw_server1.Listen("localhost", 12301)
	time.Sleep(50 * time.Millisecond)
	nw_server1.Join(address)

	ct_server2 := contact.New(kademliaid.New("1010000000000000000000000000000000000000"), "localhost:12302")
	nw_server2 := New(&ct_server2)
	go nw_server2.Listen("localhost", 12302)
	time.Sleep(50 * time.Millisecond)
	nw_server2.Join(address)

	ct_server3 := contact.New(kademliaid.New("1110000000000000000000000000000000000000"), "localhost:12303")
	nw_server3 := New(&ct_server3)
	go nw_server3.Listen("localhost", 12303)
	time.Sleep(50 * time.Millisecond)
	nw_server3.Join(address)

	ct_server4 := contact.New(kademliaid.New("0000000000000000000000000000000000000001"), "localhost:12304")
	nw_server4 := New(&ct_server4)
	go nw_server4.Listen("localhost", 12304)
	time.Sleep(50 * time.Millisecond)
	nw_server4.Join(address)
	
	target := kademliaid.New("1111000000000000000000000000000000000000")
	contacts := nw_client.LookupContact(target)

	log.Println(contacts)	


}
