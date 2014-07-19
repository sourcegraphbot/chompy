package stable

// Let's create a type to hold information about the variables in
// our program.
type NodeInfo struct {
	// What variables do we need? Probably a pointer to a type
	T Type
	// What else? We don't need the identifier name because
	// that's stored in the symbol table. There may be other
	// things but I'm not sure what they are.
	// We need to store the thing above it
	up *NodeInfo
	// We need to store the actual value.
	// TODO create a type for val
	Val interface{}
}

// Okay, let's create our Type type. Type will hold all the
// information we need to generate code for a specific type.
type Type interface {
	Equal(Type) bool
}

type Func struct {
	Args []Type
}

func (f *Func) Equal(t Type) bool {
	fn, ok := t.(*Func)
	if !ok {
		return false
	}
	if len(f.Args) != len(fn.Args) {
		return false
	}
	for i := 0; i < len(f.Args); i++ {
		if f.Args[i].Equal(fn.Args[i]) == false {
			return false
		}
	}
	return true
}

// Represents all types that are not functions
type Basic struct {
	Name string
	// this is a pointer type if true
	Pointer bool
}

func (b *Basic) Equal(t Type) bool {
	ba, ok := t.(*Basic)
	if !ok {
		return false
	}
	return b.Name == ba.Name && b.Pointer == ba.Pointer
}

type Struct struct {
	Name   string
	Fields []Type
}

func (s *Struct) Equal(t Type) bool {
	ss, ok := t.(*Struct)
	if !ok {
		return false
	}
	if s.Name != s.Name ||
		len(s.Fields) != len(ss.Fields) {
		return false
	}
	for i := 0; i < len(s.Fields); i++ {
		if s.Fields[i].Equal(ss.Fields[i]) == false {
			return false
		}
	}
	return true
}

type Stable struct {
	table  map[string]*NodeInfo
	latest *NodeInfo
	up     *Stable
}

// Creates a new Stable. It is legal to pass nil as the old Stable
func New(old *Stable) *Stable {
	return &Stable{
		table: make(map[string]*NodeInfo),
	}
}

func (s *Stable) Insert(name string, value *NodeInfo) {
	value.up = s.latest
	s.latest = value
	s.table[name] = value
}

func (s *Stable) Get(name string) (*NodeInfo, bool) {
	for tab := s; tab != nil; tab = s.up {
		n, ok := tab.table[name]
		if ok {
			return n, true
		}
	}
	return nil, false
}