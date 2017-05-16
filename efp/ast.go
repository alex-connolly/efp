package efp

import "regexp"

type field struct {
	alias string
	key   *key
	value *fieldValue
}

type key struct {
	key   string
	regex *regexp.Regexp
	min   int
	max   int
}

type fieldValue struct {
	parent   *fieldValue
	isArray  bool
	value    string
	children []*fieldValue
	min      int
	max      int
}

type element struct {
	alias                    string
	key                      *key
	parent                   *element
	parameters               []*fieldValue
	declaredTextAliases      map[string]string
	discoveredElementAliases map[string]*element
	declaredElementAliases   map[string]*element
	elements                 map[string][]*element
	fields                   map[string][]*field
	discoveredFieldAliases   map[string]*field
	declaredFieldAliases     map[string]*field
}

var standards = map[string]string{
	"string": `".*"`,
	"int":    "-[1-9]+|[0-9]+",
	"float":  "[0-9]*.[0-9]+",
	"bool":   "true|false",
	"uint":   "[0-9]+",
}

func (e *element) addStandardAliases() {
	e.declaredTextAliases = standards
}
