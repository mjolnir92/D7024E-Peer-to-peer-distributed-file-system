package kademlia

import (
	"bytes"
	"strconv"
	"time"
	"testing"
	"github.com/mjolnir92/kdfs/kademliaid"
	"github.com/mjolnir92/kdfs/contact"
)

func TestLookupContact(t *testing.T) {
	//t.Skip()
	address1 := "localhost:12400"
	ct_kademlia1 := contact.New(kademliaid.New("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), address1)
	nw_kademlia1 := New(&ct_kademlia1)
	go nw_kademlia1.Listen(address1)
	time.Sleep(50 * time.Millisecond)

	address2 := "localhost:12401"
	ct_kademlia2 := contact.New(kademliaid.New("0000000000000000000000000000000000000000"), address2)
	nw_kademlia2 := New(&ct_kademlia2)
	go nw_kademlia2.Listen(address2)
	time.Sleep(50 * time.Millisecond)
	nw_kademlia2.Join(address1)

	for i := 0; i<20; i++ {
		address := "localhost:"+strconv.Itoa(12410+i)
		ct := contact.New(kademliaid.NewRandom(), address)
		nw := New(&ct)
		go nw.Listen(address)
		time.Sleep(50*time.Millisecond)
		nw.Join(address1)
	}

	target := kademliaid.New("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	contacts := nw_kademlia1.LookupContact(target)
	for _, c := range contacts {
		if c.ID == target {
			t.Error("LookupContact did not return the correct contacts")
		}
	}
}

func TestLookupData(t *testing.T) {
	address1 := "localhost:12500"
	ct_kademlia1 := contact.New(kademliaid.New("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), address1)
	nw_kademlia1 := New(&ct_kademlia1)
	go nw_kademlia1.Listen(address1)
	time.Sleep(50 * time.Millisecond)

	address2 := "localhost:12501"
	ct_kademlia2 := contact.New(kademliaid.New("0000000000000000000000000000000000000000"), address2)
	nw_kademlia2 := New(&ct_kademlia2)
	go nw_kademlia2.Listen(address2)
	time.Sleep(50 * time.Millisecond)
	nw_kademlia2.Join(address1)
	
	for i := 0; i<20; i++ {
		address := "localhost:"+strconv.Itoa(12510+i)
		ct := contact.New(kademliaid.NewRandom(), address)
		nw := New(&ct)
		go nw.Listen(address)
		time.Sleep(50*time.Millisecond)
		nw.Join(address1)
	}

	testData := []byte("my test data")
	testData2 := []byte("should not exist")
	nw_kademlia2.KademliaStore(testData)
	time.Sleep(50 * time.Millisecond)
	data, err := nw_kademlia1.LookupData(kademliaid.NewHash(testData))
	if err != nil {
		t.Error("LookupData failed: ", err)
	} else {
		if bytes.Compare(data.GetData(), testData) != 0 {
			t.Error("LookupData failed: Wrong data returned")
		}
	}
	data, err = nw_kademlia1.LookupData(kademliaid.NewHash(testData2))
	if err == nil {
		t.Error("Requested data should not exist")
	}
}

func TestCat(t *testing.T) {
	address1 := "localhost:12600"
	ct_kademlia1 := contact.New(kademliaid.New("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), address1)
	nw_kademlia1 := New(&ct_kademlia1)
	go nw_kademlia1.Listen(address1)
	time.Sleep(50 * time.Millisecond)

	address2 := "localhost:12601"
	ct_kademlia2 := contact.New(kademliaid.New("0000000000000000000000000000000000000000"), address2)
	nw_kademlia2 := New(&ct_kademlia2)
	go nw_kademlia2.Listen(address2)
	time.Sleep(50 * time.Millisecond)
	nw_kademlia2.Join(address1)

	testData := []byte("my test data")
	nw_kademlia2.KademliaStore(testData)
	time.Sleep(50 * time.Millisecond)
	id := kademliaid.NewHash(testData)
	data := nw_kademlia1.Cat(*id)
	if bytes.Compare(data, testData) != 0 {
		t.Error("TestCat failed, wrong data")
	}
}
