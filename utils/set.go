package utils

type StringSet map[string]struct{}

func NewStringSetFromSlice(s []string) StringSet {
	return NewStringSet(s...)
}

func NewStringSet(s ...string) StringSet {
	set := make(StringSet)
	for _, e := range s {
		set.Add(e)
	}
	return set
}

func (set *StringSet) Add(e string) bool {
	_, found := (*set)[e]
	if found {
		return false
	}

	(*set)[e] = struct{}{}
	return true
}

func (set *StringSet) Contains(s ...string) bool {
	for _, e := range s {
		if _, ok := (*set)[e]; !ok {
			return false
		}
	}
	return true
}

func (set StringSet) Equals(other StringSet) bool {
	if len(set) != len(other) {
		return false
	}

	for e, _ := range set {
		if !other.Contains(e) {
			return false
		}
	}
	return true
}

func (set StringSet) Count() int {
	return len(set)
}

func (set StringSet) ToSlice() []string {
	s := make([]string, 0, set.Count())
	for e := range set {
		s = append(s, e)
	}
	return s
}
