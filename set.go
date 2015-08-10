package gpm

// Set is a simple strings set implementation
// Warning : It is not concurrency-safe
type Set map[string]struct{}

func (self Set) Has(value string) bool {
	_, exists := self[value]
	return exists
}

func (self *Set) Add(value string) {
	(*self)[value] = struct{}{}
}

func (self *Set) AddSet(set Set) {
	for value := range set {
		self.Add(value)
	}
}

func (self *Set) Remove(value string) {
	delete((*self), value)
}

func NewSet(values ...string) Set {
	s := Set(make(map[string]struct{}))
	for _, value := range values {
		s.Add(value)
	}
	return s
}
