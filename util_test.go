package efp

import (
	"testing"

	"github.com/end-r/goutil"
)

func TestStrval(t *testing.T) {
	x := "hello"
	goutil.Assert(t, x == strval("hello"), "unchanged strval failed")
	goutil.Assert(t, x == strval(`"hello"`), "changed strval failed")
}
