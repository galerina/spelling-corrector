package stringset

type StringSet map[string]struct{}

func New() *StringSet {
	set := StringSet{}
	return &set
}

func (set *StringSet) Add(s string) {
	(*set)[s] = struct{}{}
}

func (set *StringSet) Contains(s string) bool {
	_, ok := (*set)[s]
	return ok
}

func (set *StringSet) AsSlice() []string {
	slice := make([]string, len(*set))

	i := 0
	for k := range *set {
		slice[i] = k
		i++
	}

	return slice
}

func (set *StringSet) GetAny() string {
	for k := range *set {
		return k
	}

	return ""
}
