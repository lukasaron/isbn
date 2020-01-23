package isbn

import (
	"fmt"
	"testing"
)

func TestNewISBN(t *testing.T) {
	print(NewISBN("ISBN 978-0-11-000222-4"))
	print(NewISBN("0-393-04002-X"))
	print(NewISBN("9789528988886"))
	print(NewISBN("039304002X"))
	print(NewISBN("0-6198-8108-9"))
	print(NewISBN("0619881089"))
	print(NewISBN("98-7166-066-9"))
	print(NewISBN("9871660669"))
}

func print(isbn ISBN) {
	fmt.Println(isbn)
	fmt.Println(isbn.IsValid())
	isbn.Normalize()
	fmt.Println(isbn)
	fmt.Println(isbn.IsValid())
}
