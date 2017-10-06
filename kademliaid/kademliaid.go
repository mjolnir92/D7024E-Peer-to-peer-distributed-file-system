package kademliaid

import (
	"crypto/sha1"
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

func NewHash(data []byte) *T {
	newKademliaID := T{}
	newKademliaID = sha1.Sum(data)
	return &newKademliaID
}

func NewRandom() *T {
	newKademliaID := T{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

//Returns a random kademliaid with common prefix of length n to contact
func NewRandomCommonPrefix(id T, n uint8) *T {
	id_old := id
	id_rand := *NewRandom()
	//Iterate through the byte slice. Replace entire bytes in the new id if the prefix covers that byte
	//If the prefix only covers part of a byte, mask out the bytes from the prefix that should remain
	for i := 0; i < IDLength; i++ {
		if (i+1)*8 < int(n) {
			id_rand[i] = id_old[i]
		} else if (i+1)*8 == int(n) { 
			id_rand[i] = id_old[i]
			//Flip the first bit in the next byte if it wasn't the last byte
			if i != IDLength-1 {
				id_rand[i+1] = id_old[i+1] ^ (uint8(1) << 7)
			}
			break
		} else {
			//Creates a bitmask with the n%8 first bits set to 0, the rest 1
			mask := (^uint8(0) >> (n % 8))
			mask &= id_rand[i]
			//Ensure that the bit after the common prefix ends is set to 1, so that the common prefix ends there
			//Otherwise the ID would belong to a different bucket range
			if (mask &  (1 << (7 - (n % 8)))) == 0 {
				mask ^= uint8(1) << (7 - (n % 8))
			}
			mask ^= id_old[i]
			id_rand[i] = mask
			break
		}
	}
	return &id_rand
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
