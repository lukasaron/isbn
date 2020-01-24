# ISBN
ISBN library is the implementation of International Standard Book Number 
standard ISO 2108:2017. Input to this library can be a string containing 
ISBN of version 10 or version 13.  

[![GoDoc](https://godoc.org/github.com/lukasaron/isbn?status.svg)](https://godoc.org/github.com/lukasaron/isbn)
[![Build Status](https://travis-ci.com/lukasaron/isbn.svg?branch=master)](https://travis-ci.com/lukasaron/isbn)
[![Go Report Card](https://goreportcard.com/badge/github.com/lukasaron/isbn)](https://goreportcard.com/report/github.com/lukasaron/isbn)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

## Installation
```go
go get github.com/lukasaron/isbn
```

## Example of usage:
```go
package main

import (
    "fmt"
    "github.com/lukasaron/isbn"
)

func main() {
    // ISBN can be specified, otherwise the automatic detection is initiated
    book := isbn.NewISBN("9780777777770")
    fmt.Println(book, book.IsValid(), book.Version(), book.Error(), book.BarCode())

    // ISBN and ISBN-13 are the correct specification of version 13
    // the number separator could be used space or hyphen,
    // however it's not required.
    book = isbn.NewISBN("ISBN 978 0 7777 7777 0")
    fmt.Println(book, book.IsValid(), book.Version(), book.Error(), book.BarCode())

    // version 10 should be specified (by the ISO standard)
    book = isbn.NewISBN("ISBN-10 0-393-04002-X")
    fmt.Println(book, book.IsValid(), book.Version(), book.Error(), book.BarCode())

    // version 10 should be specified (by the ISO standard),
    // handles also the normalization -> conversion into version 13
    book = isbn.NewISBN("ISBN-10 039304002X")
    book.Normalize()
    fmt.Println(book, book.IsValid(), book.Version(), book.Error(), book.BarCode())

}
```