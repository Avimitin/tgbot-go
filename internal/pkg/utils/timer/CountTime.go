package timer

import "time"

type Timer struct {
	GeneratedTime float64
}

func NewTimer() *Timer {
	// New a timer instance
	now := float64(time.Now().UnixNano())
	return &Timer{now}
}

func (t *Timer) StopCounting() float64 {
	// Get duration between current time and instance generated time
	now := float64(time.Now().UnixNano())
	return now - t.GeneratedTime
}

func GetGenerateTime(t *Timer) float64 {
	// Get Unix time about when the instance generated
	return t.GeneratedTime
}
