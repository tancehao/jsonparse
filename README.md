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


* (p *Parser)Parse() (*Elem, error) {

}

