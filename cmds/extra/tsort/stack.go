package main

type stack struct {
	s []string
}

func makeStack() stack {
	return stack{}
}

func (s *stack) push(value string) {
	s.s = append(s.s, value)
}

func (s *stack) pop() string {
	if len(s.s) == 0 {
		panic("stack is empty")
	}

	result := s.s[len(s.s)-1]
	s.s = s.s[:len(s.s)-1]
	return result
}

func (s *stack) peek() string {
	if len(s.s) == 0 {
		panic("stack is empty")
	}

	return s.s[len(s.s)-1]
}
