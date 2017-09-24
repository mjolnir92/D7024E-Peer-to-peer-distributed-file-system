package network

import (
	"testing"
	"time"
	"github.com/mjolnir92/kdfs/kademliaid"
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/routingtable"
	"github.com/vmihailenco/msgpack"
)

func TestRPCs(t *testing.T) {
	//t.Skip("temp")
	// set up client
	id_client := kademliaid.New("1000000000000000000000000000000000000000")
	ct_client := contact.New(id_client, "localhost:12310")
	rt_client := routingtable.New(ct_client, 20)
	nw_client := New(5000, id_client, rt_client)
	// set up server
	id_server := kademliaid.New("0000000000000000000000000000000000000000")
	ct_server := contact.New(id_server, "localhost:12300")
	rt_server := routingtable.New(ct_server, 20)
	rt_server.AddContact(ct_client)
	nw_server := New(5000, id_server, rt_server)
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
}

func TestMarshal(t *testing.T) {
	id_client := kademliaid.New("1000000000000000000000000000000000000000")
	ct_client := contact.New(id_client, "localhost:12310")
	rt_client := routingtable.New(ct_client, 20)
	nw_client := New(5000, id_client, rt_client)
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
