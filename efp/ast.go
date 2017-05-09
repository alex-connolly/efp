package efp

type field struct {
	alias string
	key   string
	value *fieldValue
}

// fv(true, fv(false, int), fv(false, string))
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
	parent                   *element
	key                      string
	parameters               []*fieldValue
	declaredTextAliases      map[string]string
	discoveredElementAliases map[string][]*element
	declaredElementAliases   map[string][]*element
	elements                 map[string][]*element
	fields                   map[string][]*field
	discoveredFieldAliases   map[string][]*field
	declaredFieldAliases     map[string][]*field
}
