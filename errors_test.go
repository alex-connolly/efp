package efp

import (
	"testing"

	"github.com/end-r/goutil"
)

func TestOnlyOperator(t *testing.T) {
	_, errs := PrototypeString("=")
	goutil.Assert(t, errs != nil, "should be an error")
}

func TestStrayOperator(t *testing.T) {
	_, errs := PrototypeString("name :+ string")
	goutil.Assert(t, errs != nil, "should be an error")
}

func TestInvisibleAlias(t *testing.T) {
	_, errs := PrototypeString("name : boolean")
	goutil.Assert(t, errs != nil, "should be an error")
}
