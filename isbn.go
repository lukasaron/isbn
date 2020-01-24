package isbn

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

// Version specifies the ISBN version for a particular string.
type Version int

// ISBN versions.
const (
	VersionUnknown Version = 0
	Version10      Version = 10
	Version13      Version = 13
)

// ISBN default prefix.
const (
	DefaultPrefix = "978"
)

const (
	version10Mod = 11
	version13Mod = 10
)

const (
	versionXParts  = 1
	version10Parts = 4
	version13Parts = 5
)

const (
	prefixLength     = 3
	groupLength      = 5
	registrantLength = 5
	checkDigitLength = 1
)

const (
	version10GroupIdx = iota
	version10RegistrantIdx
	version10PublicationIdx
	version10CheckIdx
)

const (
	version13PrefixIdx = iota
	version13GroupIdx
	version13RegistrantIdx
	version13PublicationIdx
	version13CheckIdx
)

const headerLength = 7

var (
	isbnRegex = regexp.MustCompile(`(ISBN|isbn)?[\d\s\-]+[xX]?`)
	// only numbers or final X (for version 10 as a check digit number)
	isbnParserRegex = regexp.MustCompile(`([\dXx]+)`)
)

var (
	errWrongISBN = errors.New("wrong input ISBN format")
)

// ISBN struct defines the core ISBN logic.
type ISBN struct {
	originalISBN      string
	version           Version
	prefix            string
	registrationGroup string
	registrant        string
	publication       string
	checkDigit        string
	err               error
}

// NewISBN function creates ISBN instance based on the input string.
func NewISBN(isbnStr string) (isbn ISBN) {
	// check if the string is basic ISBN string
	match := isbnRegex.MatchString(isbnStr)
	if !match {
		isbn.err = errWrongISBN
		return isbn
	}

	numbers := isbnParserRegex.FindAllString(isbnStr, -1)
	// remove the ISBN version from numbers
	if len(numbers) > 1 && (numbers[0] == "13" || numbers[0] == "10") {
		numbers = numbers[1:]
	}

	switch len(numbers) {
	case version10Parts:
		isbn = ISBN{
			version:           Version10,
			registrationGroup: numbers[version10GroupIdx],
			registrant:        numbers[version10RegistrantIdx],
			publication:       numbers[version10PublicationIdx],
			checkDigit:        numbers[version10CheckIdx],
		}

	case version13Parts:
		isbn = ISBN{
			version:           Version13,
			prefix:            numbers[version13PrefixIdx],
			registrationGroup: numbers[version13GroupIdx],
			registrant:        numbers[version13RegistrantIdx],
			publication:       numbers[version13PublicationIdx],
			checkDigit:        numbers[version13CheckIdx],
		}

	case versionXParts:
		isbn = parseISBN(numbers[0])
	default:
		isbn.err = errWrongISBN
	}

	// set original ISBN
	isbn.originalISBN = isbnStr

	return isbn
}

// IsValid method checks the ISBN value(s) and returns true if the ISBN is valid, otherwise false.
func (isbn ISBN) IsValid() (valid bool) {
	if isbn.err != nil || len(isbn.checkDigit) != 1 || len(isbn.originalISBN) < headerLength {
		return false
	}

	originVersion := isbn.getVersionFromOriginal()
	if originVersion != VersionUnknown && originVersion != isbn.version {
		return false
	}

	switch isbn.version {
	case Version10:
		valid = isbn.calculateV10CheckDigit() == isbn.checkDigit
	case Version13:
		valid = isbn.calculateV13CheckDigit() == isbn.checkDigit
	default:
		valid = false
	}

	return valid
}

// Version method returns the current version of ISBN instance.
func (isbn ISBN) Version() Version {
	return isbn.version
}

// Normalize method converts ISBN of version 10 into version 13 and/or recalculate the check digital
// which is located at the end of this ISBN.
func (isbn *ISBN) Normalize() {
	if isbn.err != nil || isbn.version == Version13 && isbn.IsValid() {
		return
	}

	isbn.prefix = DefaultPrefix
	isbn.version = Version13
	isbn.checkDigit = isbn.calculateV13CheckDigit()
}

// Error method returns status error.
func (isbn ISBN) Error() error {
	return isbn.err
}

// String method creates a human readable format of ISBN.
func (isbn ISBN) String() string {
	switch isbn.version {
	case Version10:
		return fmt.Sprintf("ISBN-%s %s-%s-%s-%s",
			isbn.version.String(),
			isbn.registrationGroup,
			isbn.registrant,
			isbn.publication,
			isbn.checkDigit)
	case Version13:
		// we do not need to print version 13, it's implicit
		return fmt.Sprintf("ISBN %s-%s-%s-%s-%s",
			isbn.prefix,
			isbn.registrationGroup,
			isbn.registrant,
			isbn.publication,
			isbn.checkDigit)
	default:
		return ""
	}
}

// BarCode method creates an ISBN code without hyphens between each ISBN part.
func (isbn ISBN) BarCode() string {
	return fmt.Sprintf("%s%s%s%s%s",
		isbn.prefix, // version 10 has this value empty
		isbn.registrationGroup,
		isbn.registrant,
		isbn.publication,
		isbn.checkDigit)

}

// ------------------------------------------------- PRIVATE METHODS -------------------------------------------------

func (isbn ISBN) calculateV13CheckDigit() string {
	w := weightFn(isbn.version)
	sum := weightSum(isbn.prefix, w)
	sum += weightSum(isbn.registrationGroup, w)
	sum += weightSum(isbn.registrant, w)
	sum += weightSum(isbn.publication, w)

	reminder := sum % version13Mod
	if reminder == 0 {
		reminder = version13Mod
	}

	return strconv.Itoa(version13Mod - reminder)
}

func (isbn ISBN) calculateV10CheckDigit() string {
	w := weightFn(isbn.version)
	sum := weightSum(isbn.registrationGroup, w)
	sum += weightSum(isbn.registrant, w)
	sum += weightSum(isbn.publication, w)

	// reminder
	digit := version10Mod - (sum % version10Mod)
	if digit < 10 {
		return strconv.Itoa(digit)
	}

	// special case when digit == 10
	return "X"
}

func (isbn ISBN) getVersionFromOriginal() Version {
	if isbn.originalISBN[0] != 'i' && isbn.originalISBN[0] != 'I' {
		return VersionUnknown
	}

	// length of the original ISBN needs to be handled before this function is called
	possibleVersion := isbn.originalISBN[:headerLength]
	switch possibleVersion {
	case "isbn-10", "ISBN-10":
		return Version10
		// when ISBN without number -> Version 13 expected
	case "isbn", "ISBN", "isbn-13", "ISBN-13":
		return Version13
	default:
		return VersionUnknown
	}
}

// String method gets the string value of Version type
func (v Version) String() string {
	return fmt.Sprintf("%d", v)
}

// ------------------------------------------------ PRIVATE FUNCTIONS-------------------------------------------------

func parseISBN(isbnStr string) (isbn ISBN) {
	idx := 0

	// load prefix
	isbn.prefix, isbn.err = subString(isbnStr, idx, prefixLength)
	if isbn.err != nil {
		return isbn
	}

	// set versions and potentially correct prefix
	if isbn.prefix != DefaultPrefix {
		isbn.prefix = "" // version 10 doesn't have prefix
		isbn.version = Version10
	} else {
		idx += prefixLength
		isbn.version = Version13
	}

	groupLength := parseGroupLength(parseNumber(isbnStr, idx, groupLength))
	if groupLength == 0 {
		isbn.err = errWrongISBN
		return isbn
	}

	isbn.registrationGroup, isbn.err = subString(isbnStr, idx, groupLength)
	if isbn.err != nil {
		return isbn
	}

	idx += groupLength

	registrantLength := parseRegistrantLength(parseNumber(isbnStr, idx, registrantLength))
	if registrantLength == 0 {
		isbn.err = errWrongISBN
		return isbn
	}

	isbn.registrant, isbn.err = subString(isbnStr, idx, registrantLength)
	if isbn.err != nil {
		return isbn
	}

	idx += registrantLength
	lastIdx := len(isbnStr) - 1

	isbn.publication, isbn.err = subString(isbnStr, idx, lastIdx-idx)
	if isbn.err != nil {
		return isbn
	}

	isbn.checkDigit, isbn.err = subString(isbnStr, lastIdx, checkDigitLength)
	return isbn
}

func weightFn(version Version) func() int {
	switch version {
	case Version10:
		value := 10
		return func() int {
			v := value
			value--
			return v
		}
	case Version13:
		idx := -1
		values := []int{1, 3}
		return func() int {
			idx++
			return values[idx%2]

		}
	default:
		return nil
	}
}

func weightSum(number string, weight func() int) int {
	sum := 0
	for _, v := range number {
		sum += int(v-'0') * weight()
	}

	return sum
}

func parseNumber(input string, start, length int) (sum int) {
	if len(input) < start+length {
		return sum
	}

	mul := 0
	for i := start + length - 1; i >= start; i, mul = i-1, mul+1 {
		sum += int(input[i]-'0') * int(math.Pow10(mul))
	}

	return sum
}

func subString(input string, start, length int) (string, error) {
	if len(input) < start+length {
		return "", errWrongISBN
	}

	return input[start : start+length], nil
}

func parseGroupLength(group int) int {
	switch {
	case group < 60000:
		return 1
	case group < 70000:
		return 0
	case group < 80000:
		return 1
	case group < 95000:
		return 2
	case group < 99000:
		return 3
	case group < 99900:
		return 4
	case group < 99999:
		return 5
	default:
		return 0
	}
}

func parseRegistrantLength(registrant int) int {
	switch {
	case registrant < 20000:
		return 2
	case registrant < 50000:
		return 3
	case registrant < 89000:
		return 4
	case registrant < 95000:
		return 2
	case registrant < 99000:
		return 4
	case registrant < 100000:
		return 5
	default:
		return 0
	}
}
