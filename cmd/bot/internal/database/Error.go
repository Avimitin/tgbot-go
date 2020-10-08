package database

import "log"

func Pln(args ...string) {
	var text = "[DATABASE]"
	for _, arg := range args {
		text += " " + arg
	}
	log.Println(text)
}
