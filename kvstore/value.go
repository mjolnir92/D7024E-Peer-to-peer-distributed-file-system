package kvstore

import (
	"time"
)

//timestamp indicates when the key-value pair was last stored/updated?
//pin indicates whether the stored file is pinned
type Value struct {
	Timestamp time.Time
	Pin bool
	Data []byte
}


func NewValue(pin bool, data []byte) Value {
	v := Value{}
	v.Timestamp = time.Now()
	v.Pin = pin
	v.Data = data
	return v
}

func (v *Value) GetData() []byte {
	return v.Data
}

func (v *Value) GetPin() bool {
	return v.Pin
}

//Returns true if v's timestamp is earlier than u's timestamp
func (v *Value) Before(u Value) bool {
	return v.Timestamp.Before(u.Timestamp)
}