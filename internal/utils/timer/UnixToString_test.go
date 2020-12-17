package timer

import (
	"fmt"
	"testing"
	"time"
)

func TestUnixToString(t *testing.T) {
	now := time.Now().Unix()
	fmt.Println(now)
	strNow := UnixToString(now)
	fmt.Println(strNow)
}
