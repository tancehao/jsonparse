package jsonparse

const (
	T_NUMBER = iota
	T_STRING
	T_BOOL
	T_NULL
	T_ARRAY
	T_OBJECT
)

type Elem struct {
	Type     int
	Key      string //used for values in object
	Parent   *Elem
	Children map[string]*Elem

	offset int64
	limit  int64
}

//create an element, the context of the parser is needed to
//determine the element's attributes and relationships
func NewElem(t int, p *Parser, offset int64) (ele *Elem) {
	ele = &Elem{
		Type:     t,
		Key:      "",
		Parent:   p.currentContainer,
		Children: map[string]*Elem{},
		offset:   offset,
		limit:    offset,
	}

	if ele.Parent == nil {
		return
	}

	//add to it's parent's children
	if ele.Parent.Type == T_OBJECT && p.unassignedKey != "" {
		ele.Key, p.unassignedKey = p.unassignedKey, ""
		ele.Parent.Children[ret.Key] = ele
	}
	if ele.Parent.Type == T_ARRAY {
		ele.Key = strconv.Itoa(len(ele.Parent.children))
		ele.Parent.Children[ele.Key] = ele
	}

	return
}

func (ele *Elem) Content() []byte {

}

func (ele *Elem) Find(selector string) (ret *Elem) {

}

func (ele *Elem) String() string {

}

func (ele *Elem) Int64() int64 {

}

func (ele *Elem) Bool() bool {

}

func (ele *Elem) Float64() float64 {

}
