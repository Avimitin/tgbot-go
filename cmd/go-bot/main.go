package main

import (
	"log"
)

func main() {
	cfg, err := NewJsonConfig(WhereCFG(""))
	if err != nil {
		log.Fatal(err)
	}
	if err = Run(cfg); err != nil {
		log.Fatal(err)
	}
}
