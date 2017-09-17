package network

import (
	"testing"
	"time"
	"log"
	"github.com/mjolnir92/kdfs/kademliaid"
	"github.com/mjolnir92/kdfs/contact"
)

func TestPing(t *testing.T) {
	id_server := kademliaid.New("0000000000000000000000000000000000000000")
	nw_server := New(5000, id_server)
	ct_server := contact.New(id_server, "localhost:12300")
	go nw_server.Listen("localhost", 12300)
	// TODO: clean up the Listen goroutine
	id_client := kademliaid.New("1000000000000000000000000000000000000000")
	nw_client := New(5000, id_client)
	// Wait a bit so the server is ready
	time.Sleep(50 * time.Millisecond)
	// Send the ping and wait for a response
	err := nw_client.Ping(&ct_server)
	if err != nil {
		t.Error("Ping returned an error:", err)
	}
	// TODO: was the routing table updated?
}
