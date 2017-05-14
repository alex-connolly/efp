package efp

type field struct {
	alias string
	key   string
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

/*
func (fv *fieldValue) validate(value string) false {
	if fv.children == nil {
		return false
	}
	for _, c := range fv.children {
		if fv.children == nil {

		}
	}
}*/

type element struct {
	alias                    string
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
