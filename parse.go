package jsonparse

import (
	"bytes"
	"errors"
)

type Token []byte

var ErrNotJson = errors.New("unable to parse non-json data")
var ErrEOF = errors.New("EOF")

var (
	Comma        byte   = ','
	Quote        byte   = '"'
	Colon        byte   = ':'
	BraceOpen    byte   = '{'
	BraceClose   byte   = '}'
	BracketOpen  byte   = '['
	BracketClose byte   = ']'
	Null         string = "null"
	BoolTrue     string = "true"
	BoolFalse    string = "false"
)

var Separators = []byte{
	Comma,
	Colon,
	BraceOpen,
	BraceClose,
	BracketOpen,
	BracketClose,
}

var CertainValues = []string{
	Null,
	BoolTrue,
	BoolFalse,
}

type Parser struct {
	data []byte

	root             *Elem
	unassignedKey    string
	currentContainer *Elem
	stack            []byte
}

func NewParser(data []byte) (p *Parser) {
	return &Parser{data: data}
}

func (p *Parser) Parse() (root *Elem, err error) {
	var offset int64
	var ele *Elem
	for {
		token, length, err1 := readToken(p.data, offset)
		if err1 == ErrEOF {
			break
		}
		if err1 != nil {
			err = err1
			return
		}
		if length == 1 && isSeparator(token[0]) {
			tk := token[0]
			switch {
			case tk == BraceOpen || tk == BracketOpen:
				if tk == BraceOpen {
					ele = newElem(T_OBJECT, p, offset)
				} else {
					ele = newElem(T_ARRAY, p, offset)
				}
				p.currentContainer = ele
				p.stackPush(tk)
			case tk == BraceClose || tk == BracketClose:
				pre, err := p.stackPull()
				if err != nil {
					return nil, ErrNotJson
				}
				if tk == BraceClose && pre != BraceOpen {
					return nil, ErrNotJson
				}
				if tk == BracketClose && pre != BracketOpen {
					return nil, ErrNotJson
				}
				p.currentContainer.limit = offset + length
				p.currentContainer = p.currentContainer.Parent
			case tk == Comma:
				if p.currentContainer == nil || (p.currentContainer.Type != T_OBJECT && p.currentContainer.Type != T_ARRAY) {
					return nil, ErrNotJson
				}
			case tk == Colon:
				if p.currentContainer == nil || (p.currentContainer.Type != T_OBJECT && p.currentContainer.Type != T_ARRAY) {
					return nil, ErrNotJson
				}
			}
		} else {
			switch {
			case token[0] == Quote && token[length-1] == Quote:
				//string
				if p.currentContainer == nil && offset != 0 {
					return nil, ErrNotJson
				}
				//if the string is not a key in an object, create an element
				if p.currentContainer != nil && p.currentContainer.Type == T_OBJECT && p.unassignedKey == "" {
					p.unassignedKey = string(bytes.Trim(token, "\""))
				} else {
					ele = newElem(T_STRING, p, offset+1)
					ele.limit = offset + length - 1
				}
			case isCertainValue(token, length):
				if string(token) == Null {
					ele = newElem(T_NULL, p, offset)
				} else {
					ele = newElem(T_BOOL, p, offset)
				}
				ele.limit = offset + length
			default:
				ele = newElem(T_NUMBER, p, offset)
				ele.limit = offset + length
			}
		}
		if p.root == nil && ele != nil {
			for _, b := range p.data {
				ele.data = append(ele.data, b)
			}
			p.root = ele
			root = ele
		}
		offset += length
	}
	return
}

//push a token to the stack
func (p *Parser) stackPush(token byte) {
	p.stack = append(p.stack, token)
}

//pull the top element from the stack
func (p *Parser) stackPull() (byte, error) {
	length := len(p.stack)
	if length == 0 {
		return byte(0), errors.New("can't pull from an empty stack")
	}
	b := p.stack[length-1]
	p.stack = p.stack[:length-1]
	return b, nil
}

func readToken(data []byte, offset int64) (token []byte, length int64, err error) {
	if offset >= int64(len(data)) {
		return []byte{}, 0, ErrEOF
	}
	switch {
	case isSeparator(data[offset]):
		return data[offset : offset+1], 1, nil
	case data[offset] == '"':
		//string, begin with quote, keep reading until the other half comes
		for i := offset + 1; i < int64(len(data)); i++ {
			if data[i] == '"' && data[i-1] != '\\' {
				return data[offset : i+1], i - offset + 1, nil
			}
		}
		return []byte{}, 0, ErrNotJson
	default:
		//number, bool or null, read till a separator
		for i := offset + 1; i < int64(len(data)); i++ {
			if isSeparator(data[i]) && data[i-1] != '\\' {
				return data[offset:i], i - offset, nil
			}
		}
		return []byte{}, 0, ErrNotJson
	}
}

func isCertainValue(token []byte, length int64) bool {
	if length > 5 {
		return false
	}
	tk := string(token)
	if tk == BoolTrue || tk == BoolFalse || tk == Null {
		return true
	}
	return false
}

func isSeparator(b byte) bool {
	for _, s := range Separators {
		if s == b {
			return true
		}
	}
	return false
}
