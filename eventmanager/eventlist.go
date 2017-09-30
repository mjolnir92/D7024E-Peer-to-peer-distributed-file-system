package eventmanager

import (
	"sync"
	"time"
)

type eventList struct {
	List map[Event]*time.Timer
	mux sync.Mutex
}

func newEventList() *eventList {
	eventList := &eventList{}
	eventList.List = make(map[Event]*time.Timer)
	return eventList
}

func (l *eventList) insertEvent(e Event, t *time.Timer, d time.Duration) {
	l.mux.Lock()
	//If the event already exists, stop and delete the old event and insert the new.
	if val, ok := l.List[e]; ok {
		val.Stop()
		delete(l.List, e)
	}
	l.List[e] = t
	l.mux.Unlock()
}

func (l *eventList) deleteEvent(e Event) {
	l.mux.Lock()
	//If the event exists, stop the timer and delete it from the map
	if val, ok := l.List[e]; ok {
		val.Stop()
		delete(l.List, e)
	}
	l.mux.Unlock()
}

//Reset the timer of event e to the duration d if it exists. 
//This 'new' duration only last for one cycle before returning to the original duration due to the event function being a closure
func (l *eventList) resetTimer(e Event, d time.Duration) {
	l.mux.Lock()
	if val, ok := l.List[e]; ok {
		val.Reset(d)
	}
	l.mux.Unlock()
}