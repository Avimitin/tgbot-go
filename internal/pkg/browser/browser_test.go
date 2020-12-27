package browser

import "testing"

func TestBrowse(t *testing.T) {
	_, err := Browse("https://gitea.avimitin.studio")
	if err != nil {
		t.Fatal(err)
	}
}
