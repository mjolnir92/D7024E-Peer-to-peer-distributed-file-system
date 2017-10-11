package kademlia

import (
	"log"
	"bytes"
	"testing"
	"time"
	"github.com/mjolnir92/kdfs/kademliaid"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/kvstore"
	"github.com/vmihailenco/msgpack"
)

func TestRPCs(t *testing.T) {
	//t.Skip("temp")
	// set up client
	id_client := kademliaid.New("1000000000000000000000000000000000000000")
	ct_client := contact.New(id_client, "localhost:12310")
	nw_client := New(&ct_client)
	// set up server
	id_server := kademliaid.New("0000000000000000000000000000000000000000")
	ct_server := contact.New(id_server, "localhost:12300")
	nw_server := New(&ct_server)
	nw_server.routingtable.AddContact(ct_client)
	go nw_server.Listen("localhost", 12300)
	// TODO: clean up the Listen goroutine
	// Wait a bit so the server is ready
	time.Sleep(50 * time.Millisecond)
	t.Run("Ping", func(t *testing.T) {
		err := nw_client.Ping(&ct_server)
		if err != nil {
			t.Error("Ping returned an error:", err)
		}
		// TODO: was the routing table updated?
		got := nw_client.routingtable.FindClosestContacts(id_server, 1)
		if len(got) == 0 || *got[0].ID != *id_server {
			t.Error("Server was not added to routing table after Ping")
		}
	})
	t.Run("FindNode", func(t *testing.T) {
		contacts, err := nw_client.FindNode(&ct_server, id_client)
		if err != nil {
			t.Error("FindNode returned an error:", err)
		}
		if len(contacts) == 0 {
			t.Error("FindNode returned an empty contact list")
		} else if *contacts[0].ID != *ct_client.ID || contacts[0].Address != ct_client.Address {
			t.Errorf("FindNode returned an unexpected contact:\nExpected:\n%v\nGot:\n%v\n", ct_client, contacts[0])
		}
	})
	t.Run("FindValue", func(t *testing.T) {
		stored_val := kvstore.NewValue(true, []byte{255,128,0})
		id_val := kademliaid.NewHash(stored_val.GetData())
		// value is not yet stored on the server
		value, contacts, gotData, err := nw_client.FindValue(&ct_server, id_val)
		if err != nil {
			t.Error("FindValue returned an error:", err)
		}
		if gotData != false {
			t.Error("Value was found")
		}
		if len(contacts) != 0 {
			if *contacts[0].ID != *ct_client.ID || contacts[0].Address != ct_client.Address {
				t.Errorf("Value returned an unexpected contact:\nExpected:\n%v\nGot:\n%v\n", ct_client, contacts[0])
			}
		}
		// value is stored on the server
		nw_server.kvstore.Store(stored_val)
		value, contacts, gotData, err = nw_client.FindValue(&ct_server, id_val)
		if err != nil {
			t.Error("FindValue returned an error:", err)
		}
		if gotData != true {
			t.Error("Value was not found")
		}
		if !bytes.Equal(value.GetData(), stored_val.GetData()) {
			t.Errorf("Didn't get the data that was stored. Expected\n%v\nGot\n%v\n", stored_val.GetData(), value.GetData())
		}
	})
	t.Run("Store", func(t *testing.T) {
		stored_val := kvstore.NewValue(true, []byte{255,240,0})
		id_val := kademliaid.NewHash(stored_val.GetData())
		if _, ok := nw_server.kvstore.Get(*id_val); ok {
			t.Error("Test setup for Store is flawed: the value was already stored on server.")
		}
		nw_client.Store(&ct_server, &stored_val)
		// Wait so server has a chance to process the RPC
		time.Sleep(50 * time.Millisecond)
		val, ok := nw_server.kvstore.Get(*id_val)
		if !ok {
			t.Error("The key was not stored after Store RPC")
		} else if val.GetPin() != stored_val.GetPin() {
			t.Error("The stored value has the wrong pin state")
		}
	})
}

func TestMarshal(t *testing.T) {
	id_client := kademliaid.New("1000000000000000000000000000000000000000")
	ct_client := contact.New(id_client, "localhost:12310")
	nw_client := New(&ct_client)
	id_server := kademliaid.New("0000000000000000000000000000000000000000")
	expected := RPCFindNode{RPCType: FIND_NODE, Sender: *nw_client.contactMe, FindID: *id_server}
	b, err := msgpack.Marshal(&expected)
	if err != nil {
		t.Error("marshal failed")
	}
	var got RPCFindNode
	err = msgpack.Unmarshal(b, &got)
	if err != nil {
		t.Error("unmarshal failed")
	}
	log.Printf("%v", got.Sender.ID)
	// For Sender.ID, the pointers should be different, but the values should be the same
	if *got.Sender.ID != *expected.Sender.ID || got.Sender.Address != expected.Sender.Address {
		t.Errorf("Didn't get what we expected.\nexpected\n%v\ngot\n%v\n", expected, got)
	}
}
