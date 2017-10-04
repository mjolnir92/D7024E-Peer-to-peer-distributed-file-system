package eventmanager

import (
	"time"
	"github.com/mjolnir92/kdfs/kademliaid"
)

type T struct {
	list *eventList
}

func New() *T {
	return &T{list : newEventList()}
}

//Creates a new event that periodically calls the callback function f.
//The callback function can not have any arguments. If you need to send a function that takes arguments, use closures
//The eventType describes what type of event it is, it can be of any type that is comparable (no check for this at the moment)
//See the constants above for eventTypes used in the kademlia implementation
func (t *T) InsertEvent(id kademliaid.T, eventType interface{}, f func() , d time.Duration) {
	event := NewEvent(id, eventType)
	eventFunc := func() {
		f()
		//Reset the timer so that the event will periodically run. Should there be an option for non-periodic events?
		t.list.resetTimer(event, d)
	}
	timer := time.AfterFunc(d, eventFunc)
	t.list.insertEvent(event, timer, d)
}

func (t *T) DeleteEvent(id kademliaid.T, eventType interface{}) {
	event := NewEvent(id, eventType)
	t.list.deleteEvent(event)
}

func (t *T) ResetEvent(id kademliaid.T, eventType interface{}, d time.Duration) {
	event := NewEvent(id, eventType)
	t.list.resetTimer(event, d)
}