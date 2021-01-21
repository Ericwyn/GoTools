package goSet

type void struct{}

var member struct{}

type Set struct {
	privateMap map[interface{}]void
}

func NewSet() Set {
	return Set{
		privateMap: make(map[interface{}]void),
	}
}

func (set Set) Put(in interface{}) {
	set.privateMap[in] = member
}

// 删除一个
func (set Set) Remove(in interface{}) {
	delete(set.privateMap, in)
}

func (set Set) ContinueValue(in interface{}) bool {
	if _, ok := set.privateMap[in]; ok {
		return true
	}
	return false
}

func (set Set) Size() int {
	return len(set.privateMap)
}

// 删除全部
func (set Set) Clear() {
	set.privateMap = make(map[interface{}]void)
}

func (set Set) IsEmpty() bool {
	return len(set.privateMap) == 0
}
