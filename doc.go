// Package isbn presents the implementation of International Standard Book Number
// standard ISO 2108:2017. Input to this library can be string containing ISBN of version 10 or version 13.
//
// Example of usage:
//		package main
//
//		import (
//			"fmt"
//			"github.com/lukasaron/isbn"
//		)
//
//		func main() {
//			// ISBN can be specified, otherwise the automatic detection is initiated
//			book := isbn.NewISBN("9780777777770")
//			fmt.Println(book, book.IsValid(), book.Version(), book.Error(), book.BarCode())
//
//			// ISBN and ISBN-13 are the correct specification of version 13
//			// the number separator could be used space or hyphen,
//			// however it's not required.
//			book = isbn.NewISBN("ISBN 978 0 7777 7777 0")
//			fmt.Println(book, book.IsValid(), book.Version(), book.Error(), book.BarCode())
//
//			// version 10 should be specified (by the ISO standard)
//			book = isbn.NewISBN("ISBN-10 0-393-04002-X")
//			fmt.Println(book, book.IsValid(), book.Version(), book.Error(), book.BarCode())
//
//			// version 10 should be specified (by the ISO standard),
//			// handles also the normalization -> conversion into version 13
//			book = isbn.NewISBN("ISBN-10 039304002X")
//			book.Normalize()
//			fmt.Println(book, book.IsValid(), book.Version(), book.Error(), book.BarCode())
//
//		}
package isbn
