package kvstore

import (
	"time"
)

//timestamp indicates when the key-value pair was last stored/updated?
//pin indicates whether the stored file is pinned
type Value struct {
	timestamp time.Time
	pin bool
	data []byte
}


func NewValue(pin bool, data []byte) *Value {
	v := &Value{}
	v.timestamp = time.Now()
	v.pin = pin
	v.data = data
	return v
}

func (v *Value) GetData() []byte {
	return v.data
}

func (v *Value) GetPin() bool {
	return v.pin
}

//Returns true if v's timestamp is earlier than u's timestamp
func (v *Value) Before(u Value) bool {
	return v.timestamp.Before(u.timestamp)
}