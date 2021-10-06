package common

import "time"

func NewTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
