package efp

import "testing"

func TestOnlyOperator(t *testing.T) {
	_, errs := PrototypeString("=")
	assert(t, errs != nil, "should be an error")
}

func TestStrayOperator(t *testing.T) {
	_, errs := PrototypeString("name :+ string")
	assert(t, errs != nil, "should be an error")
}

func TestInvisibleAlias(t *testing.T) {
	_, errs := PrototypeString("name : boolean")
	assert(t, errs != nil, "should be an error")
}
