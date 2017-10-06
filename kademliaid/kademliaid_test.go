package kademliaid

import (
	"testing"
	"fmt"
	"strings"
)

func TestNewRandomCommonPrefix(t *testing.T) {
	id := NewRandom()
	for i := 0; i < 160; i++ {
		random := NewRandomCommonPrefix(*id, uint8(i))
		//Convert the byte slices to bit strings and compare the first i bits
		r := strings.NewReplacer(" ", "", "[", "", "]", "")
		id_str := r.Replace(fmt.Sprintf("%08b\n",*id))
		random_str := r.Replace(fmt.Sprintf("%08b\n",*random))

		if id_str[:i] != random_str[:i] {
			t.Error("TestNewRandomCommonPrefix failed, The prefix was not common")
		} 
		if id_str[i] == random_str[i] {
			fmt.Println("i: ",i)
			t.Error("TestNewRandomCommonPrefix failed, The bit after the prefix ended was still common")
		}
	}
}