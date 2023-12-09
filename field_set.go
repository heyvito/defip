package defip

type fieldSet []string

func (f fieldSet) fieldIdx(name string) int {
	for i, field := range f {
		if name == field {
			return i
		}
	}

	return -1
}
