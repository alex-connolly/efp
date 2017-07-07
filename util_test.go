package efp

import "testing"

func TestStrval(t *testing.T) {
	x := "hello"
	assert(t, x == strval("hello"), "unchanged strval failed")
	assert(t, x == strval(`"hello"`), "changed strval failed")
}
