package object

import (
	"fmt"
	"log"
	"testing"
)

func TestParseGo(t *testing.T) {
	fname := "./example/user.go"
	objects, err := ParseGo(fname)
	if err != nil {
		log.Fatalf("parse go source failed: %v", err)
	}

	fmt.Println(objects)

}
