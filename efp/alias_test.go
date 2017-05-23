package efp

import (
	"regexp"
	"testing"
)

func TestAliasingStandardAliases(t *testing.T) {
	p := createPrototypeParserString("")
	assert(t, len(p.prototype.textAliases) == len(standards), "wrong number of standard aliases")
}

func TestAliasingStandardInt(t *testing.T) {
	p := createPrototypeParserString("")
	assertNow(t, p.prototype.textAliases != nil, "text aliases is nil")
	regex, err := regexp.Compile(p.prototype.textAliases["int"].value)
	assertNow(t, err == nil, "error compiling int regex")
	assert(t, regex.MatchString("99"), "int regex didn't match")
	assert(t, regex.MatchString("0"), "int regex didn't match")
	assert(t, regex.MatchString("-100"), "int regex didn't match")

	assert(t, !regex.MatchString("aaa"), "int regex did match")
	assert(t, !regex.MatchString("-0"), "int regex did match")
}

func TestAliasingStandardUInt(t *testing.T) {
	p := createPrototypeParserString("")
	regex, err := regexp.Compile(p.prototype.textAliases["uint"].value)
	assertNow(t, err == nil, "error compiling uint regex")
	assert(t, regex.MatchString("99"), "uint regex didn't match")
	assert(t, regex.MatchString("0"), "uint regex didn't match")

	assert(t, !regex.MatchString("-100"), "uint regex did match")
	assert(t, !regex.MatchString("-0"), "uint regex did match")
}

func TestAliasingStandardFloat(t *testing.T) {
	p := createPrototypeParserString("")
	regex, err := regexp.Compile(p.prototype.textAliases["float"].value)
	assertNow(t, err == nil, "error compiling float regex")
	assert(t, regex.MatchString("99"), "float regex didn't match")
	assert(t, regex.MatchString("0"), "float regex didn't match")
	assert(t, regex.MatchString("-100"), "float regex didn't match")
	assert(t, regex.MatchString("0.5"), "float regex didn't match")

	assert(t, !regex.MatchString("aaa"), "float regex did match")
	assert(t, !regex.MatchString("-0"), "float regex did match")
}

func TestAliasingStandardBool(t *testing.T) {
	p := createPrototypeParserString("")
	regex, err := regexp.Compile(p.prototype.textAliases["bool"].value)
	assertNow(t, err == nil, "error compiling bool regex")
	assert(t, regex.MatchString("true"), "bool regex didn't match")
	assert(t, regex.MatchString("false"), "bool regex didn't match")

	assert(t, !regex.MatchString("tru"), "bool regex did match")
	assert(t, !regex.MatchString("flase"), "bool regex did match")
}

func TestAliasingFieldAlias(t *testing.T) {
	p := createPrototypeParserString(`
        alias x = name : string
        x`)
	parseFieldAlias(p)
	parseDiscoveredAlias(p)
	assertNow(t, p.errs == nil, "errs are not nil")
	assertNow(t, p.prototype.fields["name"] != nil, "name is nil")
}

func TestAliasingElementAlias(t *testing.T) {
	p := createPrototypeParserString(`
        alias x = name {

		}
        x`)
	parseElementAlias(p)
	parseDiscoveredAlias(p)
	assertNow(t, p.errs == nil, "errs are not nil")
	assertNow(t, p.prototype.fields != nil, "fields is nil")
	assertNow(t, p.prototype.fields["name"] != nil, "name is nil")
}

func TestAliasingTextAlias(t *testing.T) {
	p := createPrototypeParserString(`
        alias x = string
        name : x`)
	assertNow(t, p.prototype.fields["name"] != nil, "name is nil")
}

func TestAliasingDoubleIndirection(t *testing.T) {
	p := createPrototypeParserString(`
        alias y = string
        alias x = name : y
        x`)
	assertNow(t, p.prototype.fields["name"] != nil, "name is nil")
}

// test that element recursion is allowed
func TestAliasingRecursionValid(t *testing.T) {
	p := createPrototypeParserString(`
        alias p = x {
            p
        }
        p`)
	assertNow(t, p.prototype.elements["x"] != nil, "x is nil")
}

// test that field recursion is disallowed
func TestAliasingRecursionInvalid(t *testing.T) {
	p := createPrototypeParserString(`alias x = x`)
	assertNow(t, p.errs != nil, "errs should not be nil")
}
