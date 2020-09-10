package comm

import "time"

func NowUnix() int64 {
	return time.Now().UnixNano()
}
