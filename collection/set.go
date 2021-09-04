package collection

type Set struct {
	data map[interface{}]bool
}

func NewSet() *Set {
	return &Set{ data: map[interface{}]bool{} }
}

func (s *Set) Add(x interface{}) *Set {
	s.data[x] = true
	return s
}

func (s *Set) Del(x interface{}) *Set {
	delete(s.data, x)
	return s
}

func (s *Set) Has(x interface{}) bool {
	flag, ok := s.data[x]
	return ok && flag
}

func (s *Set) GetMap() map[interface{}]bool {
	return s.data
}