package network

import (
	"bytes"
	"testing"
	"time"
	"github.com/mjolnir92/kdfs/kademliaid"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/routingtable"
	"github.com/mjolnir92/kdfs/kvstore"
	"github.com/vmihailenco/msgpack"
)

func TestRPCs(t *testing.T) {
	//t.Skip("temp")
	// set up client
	id_client := kademliaid.New("1000000000000000000000000000000000000000")
	ct_client := contact.New(id_client, "localhost:12310")
	rt_client := routingtable.New(ct_client, 20)
	nw_client := New(5000, id_client, rt_client, nil)
	// set up server
	id_server := kademliaid.New("0000000000000000000000000000000000000000")
	ct_server := contact.New(id_server, "localhost:12300")
	rt_server := routingtable.New(ct_server, 20)
	rt_server.AddContact(ct_client)
	kvs_server := kvstore.New()
	nw_server := New(5000, id_server, rt_server, kvs_server)
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
	})
	t.Run("FindNode", func(t *testing.T) {
		contacts, err := nw_client.FindNode(&ct_server, id_client)
		if err != nil {
			t.Error("FindNode returned an error:", err)
		}
		if len(contacts) == 0 {
			t.Error("FindNode returned an empty contact list")
		}
		if *contacts[0].ID != *ct_client.ID || contacts[0].Address != ct_client.Address {
			t.Errorf("FindNode returned an unexpected contact:\nExpected:\n%v\nGot:\n%v\n", ct_client, contacts[0])
		}
	})
	t.Run("FindValue", func(t *testing.T) {
		stored_val := kvstore.NewValue(true, []byte{255,128,0})
		id_val := kademliaid.NewHash(stored_val.GetData())
		// value is not yet stored on the server
		data, contacts, gotData, err := nw_client.FindValue(&ct_server, id_val)
		if err != nil {
			t.Error("FindValue returned an error:", err)
		}
		if gotData != false {
			t.Error("Value was found")
		}
		if *contacts[0].ID != *ct_client.ID || contacts[0].Address != ct_client.Address {
			t.Errorf("Value returned an unexpected contact:\nExpected:\n%v\nGot:\n%v\n", ct_client, contacts[0])
		}
		// value is stored on the server
		nw_server.kvstore.Store(stored_val)
		data, contacts, gotData, err = nw_client.FindValue(&ct_server, id_val)
		if err != nil {
			t.Error("FindValue returned an error:", err)
		}
		if gotData != true {
			t.Error("Value was not found")
		}
		if !bytes.Equal(data, stored_val.GetData()) {
			t.Error("Didn't get the data that was stored. Expected\n%v\nGot\n%v\n", stored_val.GetData(), data)
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
	rt_client := routingtable.New(ct_client, 20)
	nw_client := New(5000, id_client, rt_client, nil)
	id_server := kademliaid.New("0000000000000000000000000000000000000000")
	expected := RPCFindNode{RPCType: FIND_NODE, SenderID: *nw_client.id, FindID: *id_server}
	b, err := msgpack.Marshal(&expected)
	if err != nil {
		t.Error("marshal failed")
	}
	var got RPCFindNode
	err = msgpack.Unmarshal(b, &got)
	if err != nil {
		t.Error("unmarshal failed")
	}
	if got != expected {
		t.Error("Didn't get what we expected")
	}
}
