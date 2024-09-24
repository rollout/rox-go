package roxx

type CoreStack struct {
	items []interface{}
}

func NewCoreStack() *CoreStack {
	return &CoreStack{}
}

func (s *CoreStack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *CoreStack) Pop() interface{} {
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *CoreStack) Peek() interface{} {
	return s.items[len(s.items)-1]
}
