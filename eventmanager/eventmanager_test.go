package eventmanager

import (
	"testing"
	"time"
	"github.com/mjolnir92/kdfs/kademliaid"
)

func TestEventManager(t *testing.T) {
	manager := New()
	id := kademliaid.NewHash([]byte("event"))
	count := 0

	f := func() {count = count + 1}
	manager.InsertEvent(*id, "EVENT",  f, 10*time.Millisecond)
	time.Sleep(35*time.Millisecond)
	if count != 3 {
		t.Error("TestEventManager failed, count was no the expected value")
	}
	manager.DeleteEvent(*id, "EVENT")
	time.Sleep(35*time.Millisecond)
	if count != 3 {
		t.Error("TestEventmanager failed, the event was not deleted")
	}
}