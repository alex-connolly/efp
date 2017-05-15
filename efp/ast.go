package efp

import "regexp"

type field struct {
	alias string
	key   string
	regex *regexp.Regexp
	value *fieldValue
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
	regex                    *regexp.Regexp
	parent                   *element
	key                      string
	parameters               []*fieldValue
	declaredTextAliases      map[string]string
	discoveredElementAliases map[string]*element
	declaredElementAliases   map[string]*element
	elements                 map[string][]*element
	fields                   map[string][]*field
	discoveredFieldAliases   map[string]*field
	declaredFieldAliases     map[string]*field
}
