package browser

import "testing"

func TestBrowse(t *testing.T) {
	Browse("https://wttr.in/zhuhai?format=v2")
}
