package jsonparse

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	T_NUMBER = iota
	T_STRING
	T_BOOL
	T_NULL
	T_ARRAY
	T_OBJECT
)

type Elem struct {
	Type        int
	Key         string //used for values in object
	Parent      *Elem
	Children    map[string]*Elem
	OrderedKeys []string //specify the order of the children

	data   []byte //used for the root element
	offset int64
	limit  int64
}

//create an element, the context of the parser is needed to
//determine the element's attributes and relationships
func newElem(t int, p *Parser, offset int64) (ele *Elem) {
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
		ele.Parent.Children[ele.Key] = ele
		ele.Parent.OrderedKeys = append(ele.Parent.OrderedKeys, ele.Key)
	}
	if ele.Parent.Type == T_ARRAY {
		ele.Key = strconv.Itoa(len(ele.Parent.Children))
		ele.Parent.Children[ele.Key] = ele
		ele.Parent.OrderedKeys = append(ele.Parent.OrderedKeys, ele.Key)
	}

	return
}

func (ele *Elem) TypeString() string {
	types := []string{"number", "string", "bool", "null", "array", "object"}
	if len(types) <= ele.Type {
		return ""
	}
	return types[ele.Type]
}

//print an element
func (ele *Elem) Print() {
	defer func() {
		if ele.Parent == nil {
			fmt.Printf("\n")
		}
	}()
	if ele.Parent == nil || ele.Parent.Type != T_OBJECT {
		printntabs(ele.level())
	}
	if ele.Type == T_STRING {
		fmt.Printf("\"%s\"", string(ele.Content()))
	} else if ele.Type == T_OBJECT || ele.Type == T_ARRAY {
		if ele.Type == T_OBJECT {
			fmt.Println("{")
			defer func() {
				printntabs(ele.level())
				fmt.Print("}")
			}()
			for _, k := range ele.OrderedKeys {
				printntabs(ele.level() + 1)
				fmt.Printf("\"%s\": ", k)
				v, _ := ele.Children[k]
				v.Print()
				fmt.Println(",")
			}
		} else {
			if ele.Parent == nil || ele.Parent.Type != T_OBJECT {
				printntabs(ele.level())
			}
			fmt.Println("[")
			defer func() {
				printntabs(ele.level())
				fmt.Print("]")
			}()
			for _, k := range ele.OrderedKeys {
				v, _ := ele.Children[k]
				v.Print()
				fmt.Println(",")
			}
		}
	} else {
		fmt.Printf("%s", string(ele.Content()))
	}
}

func printntabs(n int) {
	for i := 0; i < n; i++ {
		fmt.Printf("    ")
	}
}

//how deep the element is in the whole struct
func (ele *Elem) level() (l int) {
	for e := ele; e.Parent != nil; e = e.Parent {
		l++
	}
	return l
}

func (ele *Elem) Content() []byte {
	e := ele
	for e.Parent != nil {
		e = e.Parent
	}
	return e.data[ele.offset:ele.limit]
}

func (ele *Elem) String() string {
	return string(ele.Content())
}

func (ele *Elem) Int64() (int64, error) {
	var ErrNotInt64 error = errors.New("value is not of int64 type")
	if ele.Type != T_NUMBER {
		return 0, ErrNotInt64
	}
	i, err := strconv.ParseInt(ele.String(), 10, 64)
	if err != nil {
		return 0, ErrNotInt64
	}
	return i, nil
}

func (ele *Elem) Bool() (bool, error) {
	if ele.String() == "true" {
		return true, nil
	} else if ele.String() == "false" {
		return false, nil
	}
	return false, errors.New("value is not a boolean one")
}

func (ele *Elem) Float64() (float64, error) {
	var ErrNotFloat64 error = errors.New("value is not of float64 type")
	i, err := strconv.ParseFloat(ele.String(), 64)
	if err != nil {
		return 0.0, ErrNotFloat64
	}
	return i, nil
}

func (ele *Elem) Slice() ([]string, error) {
	var ErrNotSlice error = errors.New("value is not of slice type")
	if ele.Type != T_ARRAY && ele.Type != T_OBJECT {
		return []string{}, ErrNotSlice
	}
	s := make([]string, len(ele.Children))
	for i, key := range ele.OrderedKeys {
		s[i] = ele.Children[key].String()
	}
	return s, nil
}

func (ele *Elem) Map() (map[string]string, error) {
	var ErrNotMap error = errors.New("value is not of map type")
	if ele.Type != T_OBJECT {
		return map[string]string{}, ErrNotMap
	}
	s := make(map[string]string, len(ele.Children))
	for _, key := range ele.OrderedKeys {
		s[key] = ele.Children[key].String()
	}
	return s, nil
}

func (ele *Elem) Find(path string) (ret *Elem, err error) {
	if path == "" {
		return ele, nil
	}
	switch {
	case ele.Type == T_NUMBER || ele.Type == T_STRING || ele.Type == T_BOOL || ele.Type == T_NULL:
		return nil, fmt.Errorf("simple values are not extractable")
	case ele.Type == T_OBJECT || ele.Type == T_ARRAY:
		selector, err := readSelector(path)
		if err != nil {
			return nil, err
		}
		//TODO need validation
		var key string
		if ele.Type == T_OBJECT {
			key = selector[1:]
		} else {
			key = selector[1 : len(selector)-1]
		}

		newEle, ok := ele.Children[key]
		if !ok {
			return nil, fmt.Errorf("key or index not exist: %s", key)
		}
		return newEle.Find(path[len(selector):])
	}
	return nil, fmt.Errorf("non-json element")
}

func (e *Elem) Select(path ...string) (values []string) {
	for _, p := range path {
		var value string
		ele, err := e.Find(p)
		if err == nil {
			value = ele.String()
		}
		values = append(values, value)
	}
	return
}

//range the children with order
//func (ele *Elem) MapChildren(f func(*Elem)) {
func (ele *Elem) IterateChildren(f func(*Elem)) {
    for _, key := range ele.OrderedKeys {
		f(ele.Children[key])
	}
}

//read a selector from a string
//example: readSelector(".key1[index1].key2")  returns key1, nil
func readSelector(path string) (string, error) {
	if path == "" || (path[0] != '.' && path[0] != '[') {
		return "", fmt.Errorf("%s is not a meaningful selector which should begin with . or [", path)
	}
	if path[0] == '.' {
		for i := 1; i < len(path); i++ {
			if path[i] == '.' || path[i] == '[' {
				return path[:i], nil
			}
		}
		return path, nil
	} else {
		for i := 1; i < len(path); i++ {
			if path[i] == ']' {
				return path[:i+1], nil
			}
		}
		return "", fmt.Errorf("%s is not a meaningful selector which should begin with . or [", path)
	}
}

func Get(json []byte, path string) (r *Elem, err error) {
	parser := NewParser(json)
	root, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return root.Find(path)
}
