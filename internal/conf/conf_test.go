package conf

import "testing"

func TestInGroup(t *testing.T) {
	c := Config{
		Groups: []int64{
			123456,
			789101,
			649191333,
			784143491,
			4839417104,
		},
	}

	if !c.InGroups(649191333) {
		t.Errorf("Expect given id %d in groups but got false", 649191333)
		return
	}

	if c.InGroups(293194189) {
		t.Errorf("Unexpected given id in groups")
	}
}
