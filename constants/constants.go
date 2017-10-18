package constants

import (
	"time"
)

const (
	ALPHA = 3
	K = 20

	TIMEOUT = 500 * time.Millisecond

	PUBLISH_TIME = 24 * time.Hour
	REPUBLISH_TIME = time.Hour
	EXPIRE_TIME = 24 * time.Hour
	BUCKET_REFRESH = time.Hour

	PUBLISH = "PUBLISH"
	REPUBLISH = "REPUBLISH"
	EXPIRE = "EXPIRE"
)
