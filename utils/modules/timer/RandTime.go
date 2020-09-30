package timer

import (
	"math/rand"
	"time"
)

func AddRandTimeFromNow() int64 {
	randTime := rand.Int63n(300)
	return time.Now().Unix() + randTime
}
