package constants

import (
	"time"
)

const (
	ALPHA = 3
	K = 20

	TIMEOUT = time.Duration(5000)

	PUBLISH_TIME = 24 * time.Hour
	REPUBLISH_TIME = time.Hour
	EXPIRE_TIME = 24 * time.Hour
	BUCKET_REFRESH = time.Hour

	PUBLISH = "PUBLISH"
	REPUBLISH = "REPUBLISH"
	EXPIRE = "EXPIRE"
)

//eventTypes for events in the kademlia implementation
//bucket refresh events should just have an int to signify bucket index