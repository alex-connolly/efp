package efp

import (
	"fmt"
	"testing"
)

func testFile(name string) string {
	return fmt.Sprintf("tests/%s", name)
}

func TestLargeFile(t *testing.T) {
	p, errs := PrototypeFile(testFile("large.efp"))
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p != nil, "p should not be nil")
}

func TestLargeFiles(t *testing.T) {
	p, errs := PrototypeFile(testFile("large.efp"))
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p != nil, "p should not be nil")
}
