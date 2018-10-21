package config

import (
	"log"
	"testing"
)

func TestLoadFile(t *testing.T) {
	f := "test.yaml"
	c, b, err := LoadFile(f)
	log.Println(err)
	log.Println(c)
	log.Println(string(b))
}
