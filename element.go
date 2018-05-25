package jsonparse

import (
    "strconv"
    "fmt"
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
	Type     int
	Key      string //used for values in object
	Parent   *Elem
	Children map[string]*Elem
    childrenKeys []string    //specify the order of the children

	data []byte  //used for the root element
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
        ele.Parent.Children[ele.Key] = ele
	    ele.Parent.childrenKeys = append(ele.Parent.childrenKeys, ele.Key)
    }
	if ele.Parent.Type == T_ARRAY {
		ele.Key = strconv.Itoa(len(ele.Parent.Children))
		ele.Parent.Children[ele.Key] = ele
	    ele.Parent.childrenKeys = append(ele.Parent.childrenKeys, ele.Key)
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
    e := ele
    for e.Parent != nil {
        e = e.Parent
        //fmt.Print("  ")
    }
    if ele.Type != T_OBJECT && ele.Type != T_ARRAY {
        fmt.Print(string(ele.Content()))
    } else {
        if ele.Type == T_OBJECT {
            fmt.Println("{")
            defer fmt.Print("}")
            for _, k := range ele.childrenKeys {
                fmt.Printf("%s: ", k)
                v, _ := ele.Children[k]
                v.Print()
                fmt.Println(",")
            }
        } else {
            fmt.Println("[")
            defer fmt.Println("]")
            for _, k := range ele.childrenKeys {
                v, _ := ele.Children[k]
                v.Print()
                fmt.Println(",")
            }
        }
    }
}

func (ele *Elem) Content() []byte {
    e := ele
    for e.Parent != nil {
        e = e.Parent
    }
    return e.data[ele.offset:ele.limit]
}

func (ele *Elem) Find(selector string) (ret *Elem) {
    return nil
}

func (ele *Elem) String() string {
    return ""
}

func (ele *Elem) Int64() int64 {
    return 0
}

func (ele *Elem) Bool() bool {
    return false
}

func (ele *Elem) Float64() float64 {
    return 0.0
}
