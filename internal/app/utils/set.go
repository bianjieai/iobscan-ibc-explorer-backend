package utils

type StringSet map[string]struct{}

func NewStringSetFromStr(str ...string) StringSet {
	set := NewStringSet()
	set.AddAll(str...)
	return set
}

func NewStringSet() StringSet {
	return make(map[string]struct{})
}

func (set StringSet) Add(str string) {
	set[str] = struct{}{}
}

func (set StringSet) AddAll(str ...string) {
	for _, v := range str {
		set[v] = struct{}{}
	}
}

func (set StringSet) Remove(str string) {
	delete(set, str)
}

func (set StringSet) RemoveAll(str ...string) {
	for _, v := range str {
		delete(set, v)
	}
}

func (set StringSet) ToSlice() (res []string) {
	for k := range set {
		res = append(res, k)
	}
	return
}
