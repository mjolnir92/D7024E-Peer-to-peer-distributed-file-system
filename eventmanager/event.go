package eventmanager

import (
	"github.com/mjolnir92/kdfs/kademliaid"
)

type Event struct {
	Id kademliaid.T
	EventType interface{}
}

func NewEvent(id kademliaid.T, i interface{}) Event {
	return Event{Id: id, EventType: i}
}