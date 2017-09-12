package kademliaid

import (
	"encoding/hex"
	"math/rand"
)

const IDLength = 20

type T [IDLength]byte

func New(data string) *T {
	decoded, _ := hex.DecodeString(data)

	newKademliaID := T{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i]
	}

	return &newKademliaID
}

func NewRandom() *T {
	newKademliaID := T{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

func (kademliaID T) Less(otherKademliaID *T) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

func (kademliaID T) Equals(otherKademliaID *T) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

func (kademliaID T) CalcDistance(target *T) *T {
	result := T{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

func (kademliaID *T) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}
