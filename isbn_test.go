package isbn

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewISBN_ISBN13(t *testing.T) {
	isbn := NewISBN("ISBN 978-0-7777-7777-0")
	if isbn.err != nil {
		t.Errorf("error not expected, got: %v", isbn.err)
	}

	if !reflect.DeepEqual(isbn, isbn13) {
		t.Errorf("expected: %+v, got: %+v", isbn13, isbn)
	}

	isbn = NewISBN("ISBN 978 0 7777 7777 0")
	if isbn.err != nil {
		t.Errorf("error not expected, got: %v", isbn.err)
	}

	if !reflect.DeepEqual(isbn, isbn13) {
		t.Errorf("expected: %+v, got: %+v", isbn13, isbn)
	}

	isbn = NewISBN("978-0-7777-7777-0")
	if isbn.err != nil {
		t.Errorf("error not expected, got: %v", isbn.err)
	}

	if !reflect.DeepEqual(isbn, isbn13) {
		t.Errorf("expected: %+v, got: %+v", isbn13, isbn)
	}

	isbn = NewISBN("978 0 7777 7777 0")
	if isbn.err != nil {
		t.Errorf("error not expected, got: %v", isbn.err)
	}

	if !reflect.DeepEqual(isbn, isbn13) {
		t.Errorf("expected: %+v, got: %+v", isbn13, isbn)
	}

	isbn = NewISBN("ISBN 9780777777770")
	if isbn.err != nil {
		t.Errorf("error not expected, got: %v", isbn.err)
	}

	if !reflect.DeepEqual(isbn, isbn13) {
		t.Errorf("expected: %+v, got: %+v", isbn13, isbn)
	}

	isbn = NewISBN("9780777777770")
	if isbn.err != nil {
		t.Errorf("error not expected, got: %v", isbn.err)
	}

	if !reflect.DeepEqual(isbn, isbn13) {
		t.Errorf("expected: %+v, got: %+v", isbn13, isbn)
	}
}

func TestNewISBN_ISBN10(t *testing.T) {
	isbn := NewISBN("ISBN 0-393-04002-X")
	if isbn.err != nil {
		t.Errorf("error not expected, got: %v", isbn.err)
	}

	if !reflect.DeepEqual(isbn, isbn10) {
		t.Errorf("expected: %+v, got: %+v", isbn13, isbn)
	}

	isbn = NewISBN("ISBN 0 393 04002 X")
	if isbn.err != nil {
		t.Errorf("error not expected, got: %v", isbn.err)
	}

	if !reflect.DeepEqual(isbn, isbn10) {
		t.Errorf("expected: %+v, got: %+v", isbn10, isbn)
	}

	isbn = NewISBN("0-393-04002-X")
	if isbn.err != nil {
		t.Errorf("error not expected, got: %v", isbn.err)
	}

	if !reflect.DeepEqual(isbn, isbn10) {
		t.Errorf("expected: %+v, got: %+v", isbn10, isbn)
	}

	isbn = NewISBN("0 393 04002 X")
	if isbn.err != nil {
		t.Errorf("error not expected, got: %v", isbn.err)
	}

	if !reflect.DeepEqual(isbn, isbn10) {
		t.Errorf("expected: %+v, got: %+v", isbn10, isbn)
	}

	isbn = NewISBN("039304002X")
	if isbn.err != nil {
		t.Errorf("error not expected, got: %v", isbn.err)
	}

	if !reflect.DeepEqual(isbn, isbn10) {
		t.Errorf("expected: %+v, got: %+v", isbn10, isbn)
	}
}

func TestISBN_IsValid(t *testing.T) {
	if !isbn10.IsValid() || !isbn13.IsValid() {
		t.Error("should be valid")
	}

	isbn1 := ISBN{
		version:           Version10,
		prefix:            "",
		registrationGroup: "0",
		registrant:        "393",
		publication:       "04002",
		checkDigit:        "2",
	}

	isbn2 := ISBN{
		version:           Version13,
		prefix:            "978",
		registrationGroup: "0",
		registrant:        "7777",
		publication:       "7777",
		checkDigit:        "8",
	}

	if isbn1.IsValid() || isbn2.IsValid() {
		t.Error("check digit is not correct")
	}

	isbn1.err = errors.New("something happen")
	if isbn1.IsValid() {
		t.Error("ISBN should have status error thus it's not valid")
	}
}

func TestISBN_Error(t *testing.T) {
	if isbn10.Error() != nil {
		t.Errorf("error not expected, got: %v", isbn10.Error())
	}

	if isbn13.Error() != nil {
		t.Errorf("error not expected, got: %v", isbn13.Error())
	}

	isbn := ISBN{
		version:           Version10,
		prefix:            "",
		registrationGroup: "0",
		registrant:        "393",
		publication:       "04002",
		checkDigit:        "2",
		err:               errors.New("something happen"),
	}

	if isbn.Error() == nil {
		t.Error("expected error, got nothing")
	}
}

func TestISBN_Version_13(t *testing.T) {
	isbn := NewISBN("ISBN 978-0-7777-7777-0")
	if isbn.Version() != Version13 {
		t.Error("wrong version detected")
	}

	isbn = NewISBN("ISBN 978 0 7777 7777 0")
	if isbn.Version() != Version13 {
		t.Error("wrong version detected")
	}

	isbn = NewISBN("978-0-7777-7777-0")
	if isbn.Version() != Version13 {
		t.Error("wrong version detected")
	}

	isbn = NewISBN("978 0 7777 7777 0")
	if isbn.Version() != Version13 {
		t.Error("wrong version detected")
	}

	isbn = NewISBN("ISBN 9780777777770")
	if isbn.Version() != Version13 {
		t.Error("wrong version detected")
	}

	isbn = NewISBN("9780777777770")
	if isbn.Version() != Version13 {
		t.Error("wrong version detected")
	}
}

func TestISBN_Version_10(t *testing.T) {
	isbn := NewISBN("ISBN 0-393-04002-X")
	if isbn.Version() != Version10 {
		t.Error("wrong version detected")
	}

	isbn = NewISBN("ISBN 0 393 04002 X")
	if isbn.Version() != Version10 {
		t.Error("wrong version detected")
	}

	isbn = NewISBN("0-393-04002-X")
	if isbn.Version() != Version10 {
		t.Error("wrong version detected")
	}

	isbn = NewISBN("0 393 04002 X")
	if isbn.Version() != Version10 {
		t.Error("wrong version detected")
	}

	isbn = NewISBN("039304002X")
	if isbn.Version() != Version10 {
		t.Error("wrong version detected")
	}
}

func TestISBN_Normalize(t *testing.T) {
	isbn1 := ISBN{
		version:           Version10,
		prefix:            "",
		registrationGroup: "0",
		registrant:        "393",
		publication:       "04002",
		checkDigit:        "X",
	}

	isbn2 := ISBN{
		version:           Version13,
		prefix:            "978",
		registrationGroup: "0",
		registrant:        "7777",
		publication:       "7777",
		checkDigit:        "0",
	}

	isbn2.Normalize()
	if !reflect.DeepEqual(isbn2, isbn13) {
		t.Error("normalization of correct ISBN version 13 has no effect")
	}

	isbn1.Normalize()
	if isbn1.String() != "ISBN 978-0-393-04002-9" {
		t.Errorf("ISBN version 10 should be normalized into: %v, got %v",
			"ISBN 978-0-393-04002-5", isbn1.String())
	}

	isbn2.checkDigit = "5"
	isbn2.Normalize()
	if isbn2.checkDigit != isbn13.checkDigit {
		t.Errorf("normalization should recalculate check digit, expected: %v, got %v",
			0, isbn2.checkDigit)
	}
}

func TestISBN_String(t *testing.T) {
	if isbn10.String() != "ISBN 0-393-04002-X" {
		t.Errorf("expected: %v, got: %v", "ISBN 0-393-04002-X", isbn10.String())
	}

	if isbn13.String() != "ISBN 978-0-7777-7777-0" {
		t.Errorf("expected: %v, got: %v", "ISBN 978-0-7777-7777-0", isbn13.String())
	}
}

func TestISBN_BarCode(t *testing.T) {
	if isbn10.BarCode() != "039304002X" {
		t.Errorf("expected: %v, got: %v", "039304002X", isbn10.BarCode())
	}

	if isbn13.BarCode() != "9780777777770" {
		t.Errorf("expected: %v, got: %v", "9780777777770", isbn13.BarCode())
	}
}

// ------------------------------------------------------ DATA ------------------------------------------------------

var isbn10 = ISBN{
	version:           Version10,
	prefix:            "",
	registrationGroup: "0",
	registrant:        "393",
	publication:       "04002",
	checkDigit:        "X",
}

var isbn13 = ISBN{
	version:           Version13,
	prefix:            "978",
	registrationGroup: "0",
	registrant:        "7777",
	publication:       "7777",
	checkDigit:        "0",
}
