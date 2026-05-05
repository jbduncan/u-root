package main

func newIDAssigner() *idAssigner {
	return &idAssigner{
		valueToID: make(map[string]int, 256),
		idToValue: make([]string, 0, 256),
	}
}

type idAssigner struct {
	valueToID map[string]int
	idToValue []string
}

func (i *idAssigner) assignID(value string) int {
	if id, ok := i.valueToID[value]; ok {
		return id
	}

	id := len(i.idToValue)
	i.idToValue = append(i.idToValue, value)
	i.valueToID[value] = id
	return id
}

func (i *idAssigner) valueFor(id int) string {
	return i.idToValue[id]
}
