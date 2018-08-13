# jsonparse
a friendly tool for extracting the values from a json

---

* [Install](#install)
* [Api](#api)

---

## Install

```sh
go get -u github.com/tancehao/jsonparse
```

## api
### Parser
```go
    type Parser struct {
        //unexported fields
    }
```

* NewParser(data []byte) *Parser

    Initiate a parser with the original json bytes.


* (p *Parser)Parse() (*Elem, error)

    Parse the original bytes into a structure with whom one can friendly work with.


### Elem
```go
    type Elem struct {
        //the type of a element, it can be one of T_NUMBER, T_STRING, T_BOOL, T_NULL, T_ARRAY or T_OBJECT
        Type        int
        //the key of this element in its parent. It's used for children of objects and arrays
        Key         string
        //the parent of an element
        Parent      *Elem
        //the children of an element, it may not empty when it's an object or an array
        Children    map[string]*Elem
        //the ordered keys of its children
        OrderedKeys []string
    
        //unexported fields
    }
```


* (ele *Elem)TypeString() string 

    Return the type of an element in string format, it can be one of "number", "string", "bool", ""null", "array", "object"


* (ele *Elem)Print() 

    Print the element with level-dependent indents.


* (ele *Elem)Content() []byte

    Return the original bytes of an element


* (ele *Elem)String() string

    Implements fmt.Stringer


* (ele *Elem)Int64() (int64, error)

    Convert the element into an int64 if it's of type T_NUMBER, error was returned otherwise.


* (ele *Elem)Bool() (bool, error)
    
    Convert the element into an bool if it's of type T_BOOL, error was returned otherwise.


* (ele *Elem)Float64() (float64, error)
    
    Convert the element into an float64 if it's a float, error was returned otherwise.


* (ele *Elem)Slice() ([]string, error)

    Return the children of an object or an array in a list of strings. Each element in the children was stringified. And the string keys where replaced with int indexes.
    #### Example
    ```go
        json := []byte(`{"foo1":"bar1","foo2":"bar2","foo3":"bar3"}`)
        data, _ := jsonparse.NewParser(json).Parse()
        fmt.Println(data.Slice())
        // ["bar1", "bar2", "bar3"]
    ```


* (ele *Elem)Map() (map[string]string, error)

    Smilar to Slice(), except that the children have string keys.



* (ele *Elem)Find(path string) (*Elem, error) 

    Find an element in an bigger one specified by a path. This method make one be able to extract the values from a json the way they did in some dynamic languages.
    #### Example
    ```go
        json := []byte(`
            {"success": true, 
             "data": [
               {"foo1":"bar1"},
               {"foo2":"bar2"}
             ]
            }`)
        data, _ := jsonparse.NewParser(json).Parse()
        fmt.Println(data.Find(".data[2].foo2"))
        //bar2
    ```


* (ele *Elem)Select(path ...string) []string

    Get more than one value from a json once. Each  value can be specified with a path, and found at the result list the order it was apecified.


* (ele *Elem)IterateChildren(f func(*Elem))

    Apply a function to each child of an element.
