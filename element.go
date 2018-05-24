package jsonparse

const (
	T_NUMBER = iota
	T_STRING
	T_BOOL
	T_ARRAY
	T_OBJECT
)

type Elem struct {
	Type     int
	Key      string //used for values in object
	Parent   *Elem
	Children []*Elem
	Map      map[string]*Elem

	offset int64
	limit  int64
}

func NewElem(t int, parent *Elem, offset int64) *Elem {
	return &Elem{
		Type:     t,
		Key:      "",
		Parent:   parent,
		Children: []*Elem{},
		Map:      map[string]*Elem{},
		offset:   offset,
		limit:    offset,
	}
}

//add a element to it's parent's children
func (ele *Elem) joinSiblings(key string) {
	if ele.Parent {
		if ele.Parent.Type == T_OBJECT {
			ele.Parent.Map[key] = ele
			ele.Key = key
		} else {
			ele.Parent.Children = append(ele.Parent.Children, ele)
		}
	}
}
