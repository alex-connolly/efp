package efp

import (
	"fmt"
	"testing"
)

func testFile(name string) string {
	return fmt.Sprintf("test_files/%s", name)
}

func TestUnknownFile(t *testing.T) {
	_, errs := PrototypeFile("not_found.efp")
	assert(t, errs != nil, "errs should not be nil")
}

func TestEmptyFile(t *testing.T) {
	p, errs := PrototypeFile(testFile("empty.efp"))
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p != nil, "p should not be nil")
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
	//e, errs := p.ValidateFiles("")
}

func TestVMGenFile(t *testing.T) {
	p, errs := PrototypeFile(testFile("vm.efp"))
	assert(t, errs == nil, "errs should be nil")
	assertNow(t, p != nil, "p should not be nil")
}
