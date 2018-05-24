package jsonparse

import (
	"errors"
)

type Token []byte

var ErrNotJson = errors.New("unable to parse non-json data")
var ErrEOF = errors.New("EOF")

var (
	Comma        = ','
	Quote        = '"'
	Colon        = ':'
	BraceOpen    = '{'
	BraceClose   = '}'
	BracketOpen  = '['
	BracketClose = ']'
	Null         = "null"
	BoolTrue     = "true"
	BoolFalse    = "false"
)

var Separators = []byte{
	Comma,
	Quote,
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

    unassignedKey string
    currentElem *Elem
	currentContainer *Elem
    stack []byte
}

func NewParser(data []byte) (p *Parser) {
	return &Parser{data: data}
}

func (p *Parser) Parse() (err error) {
	var offset int64
	for {
		token, length, err := ReadToken(p.data, offset)
		if err == ErrEOF {
			break
		}
		if err != nil {
			return
		}
		if length == 1 && IsSeparator(token[0]) {
            tk := token[0]
            switch {
            case tk == BraceOpen || tk == BracketOpen {
                var ele *Elem
                if tk == BraceOpen {
                    ele = NewElem(T_OBJECT, p.currentContainer, offset)
                } else {
                    ele = NewElem(T_ARRAY, p.currentContainer, offset)
                }
                ele.joinSiblings(p.unassignedKey)
                p.currentContainer = ele
                p.stackPush(tk)
            }
            case tk == BraceClose || tk == BracketClose:
                pre, err := p.stackPull()
                if err != nil{
                    return ErrNotJson
                }
                if tk == BraceClose && pre != BraceClose {
                    return ErrNotJson
                }
                if tk == BracketClose && pre != BracketOpen {
                    return ErrNotJson
                }
                p.currentContainer.limit = offset + length
                p.currentContainer = o.currentContainer.Parent
            case tk == Comma:
                if p.currentContainer == nil || (p.currentContainer.Type != T_OBJECT && p.currentContainer.Type != T_ARRAY) {
                    return ErrNotJson
                }
                p.currentElem.limit = offset
                p.currentElem = nil
            case tk == Colon:
                if p.currentContainer == nil || (p.currentContainer.Type != T_OBJECT && p.currentContainer.Type != T_ARRAY) {
                    return ErrNotJson
                }
            }
        } else {
            switch {
            case token[0] == Quote && token[length-1] == Quote {
                if p.currentContainer == nil && offset != 0 {
                    return ErrNotJson
                }
                if offset == 0 { // the json is a single string
                    ele = NewElem(T_STRING, nil, 0)
                    ele.limit = length
                }
                if p.currentContainer != nil {
                    
                }
            }
            }
        }
        offset += length
	}
	return nil
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

func ReadToken(data []byte, offset int64) (token []byte, length int64, err error) {
	if offset >= int64(len(data)) {
		return []byte{}, 0, ErrEOF
	}
	switch {
	case IsSeparator(data[offset])
        return data[offset:offset+1], 1 nil
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
			if IsSeparator(data[i]) && data[i-1] != '\\' {
				return data[offset:i], i - offset, nil
			}
		}
		return []byte{}, 0, ErrNotJson
	}
}

func IsSeparator(b byte) bool {
	for _, s := range Separators {
		if s == b {
			return true
		}
	}
	return false
}