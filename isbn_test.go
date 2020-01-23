package isbn

import (
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
