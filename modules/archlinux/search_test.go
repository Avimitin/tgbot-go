package archlinux

import (
  "testing"
)

func TestSearch(t *testing.T) {
  resp, err := SearchAll("fcitx5")
  if err != nil {
    t.Fatal(err)
  }

  if len(resp.Results) < 1 {
    t.Fatalf("No result")
  }

  if resp.Results[0].Pkgname != "fcitx5" {
    t.Fatalf("Expect pkgname == fcitx5, got %s", resp.Results[0].Pkgname)
  }
}
