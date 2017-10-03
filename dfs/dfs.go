package dfs

import (
	"github.com/mjolnir92/kdfs/contact"
	"github.com/mjolnir92/kdfs/kademliaid"
)

//Interface describing a distributed file system.
type T interface {
	LookupContact(kademliaid.T) []contact.T
}