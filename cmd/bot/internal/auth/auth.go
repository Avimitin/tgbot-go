package auth

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var creator int

func Init(filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(file, &m)
	fmt.Printf("m: %v", m)
	return nil
}

func IsCreator(uid int) bool {
	return uid == creator
}
