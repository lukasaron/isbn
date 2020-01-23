package isbn

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

type Version int

const (
	Version10 Version = 10
	Version13 Version = 13
)

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
	prefixLength      = 3
	groupLength       = 5
	registrantLength  = 5
	publicationLength = 6
	checkDigitLength  = 1
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

var (
	// only numbers or final X (for version 10 as a check digit number)
	isbnRegex = regexp.MustCompile(`(\d*X?x?)`)
)

var (
	errWrongISBN = errors.New("wrong input ISBN format")
)

type ISBN struct {
	version           Version
	prefix            string
	registrationGroup string
	registrant        string
	publication       string
	checkDigit        string
	err               error
}

func NewISBN(isbn string) ISBN {
	return parseISBN(isbn)

}

func (isbn ISBN) IsValid() bool {
	if isbn.err != nil || len(isbn.checkDigit) != 1 {
		return false
	}

	switch isbn.version {
	case Version10:
		return isbn.calculateV10CheckDigit() == isbn.checkDigit
	case Version13:
		return isbn.calculateV13CheckDigit() == isbn.checkDigit
	default:
		return false
	}
}

func (isbn ISBN) Version() Version {
	return isbn.version
}

func (isbn *ISBN) Normalize() {
	if isbn.err != nil || isbn.version == Version13 && isbn.IsValid() {
		return
	}

	isbn.prefix = DefaultPrefix
	isbn.version = Version13
	isbn.checkDigit = isbn.calculateV13CheckDigit()
}

func (isbn ISBN) Error() error {
	return isbn.err
}

func (isbn ISBN) String() string {
	switch isbn.version {
	case Version10:
		return fmt.Sprintf("ISBN %s-%s-%s-%s",
			isbn.registrationGroup,
			isbn.registrant,
			isbn.publication,
			isbn.checkDigit)
	case Version13:
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

	return "X"
}

// ------------------------------------------------ PRIVATE FUNCTIONS-------------------------------------------------

func parseISBN(isbn string) ISBN {
	numbers := isbnRegex.FindAllString(isbn, -1)
	switch len(numbers) {
	case version10Parts:
		return ISBN{
			version:           Version10,
			registrationGroup: numbers[version10GroupIdx],
			registrant:        numbers[version10RegistrantIdx],
			publication:       numbers[version10PublicationIdx],
			checkDigit:        numbers[version10CheckIdx],
		}

	case version13Parts:
		return ISBN{
			version:           Version13,
			prefix:            numbers[version13PrefixIdx],
			registrationGroup: numbers[version13GroupIdx],
			registrant:        numbers[version13RegistrantIdx],
			publication:       numbers[version13PublicationIdx],
			checkDigit:        numbers[version13CheckIdx],
		}

	case versionXParts:
		return parseLineISBN(numbers[0])

	default:
		return ISBN{err: errWrongISBN}
	}
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

func parseLineISBN(line string) (isbn ISBN) {
	start := 0

	// load prefix
	isbn.prefix, isbn.err = subString(line, start, prefixLength)
	if isbn.err != nil {
		return isbn
	}

	// set versions and potentially correct prefix
	if isbn.prefix != DefaultPrefix {
		isbn.prefix = ""
		isbn.version = Version10
	} else {
		start += 3
		isbn.version = Version13
	}

	groupLength := parseGroupLength(parseNumber(line, start, groupLength))
	if groupLength == 0 {
		isbn.err = errWrongISBN
		return isbn
	}

	isbn.registrationGroup, isbn.err = subString(line, start, groupLength)
	if isbn.err != nil {
		return isbn
	}

	start += groupLength

	registrantLength := parseRegistrantLength(parseNumber(line, start, registrantLength))
	if registrantLength == 0 {
		isbn.err = errWrongISBN
		return isbn
	}

	isbn.registrant, isbn.err = subString(line, start, registrantLength)
	if isbn.err != nil {
		return isbn
	}

	start += registrantLength

	isbn.publication, isbn.err = subString(line, start, publicationLength-registrantLength)
	if isbn.err != nil {
		return isbn
	}

	isbn.checkDigit, isbn.err = subString(line, len(line)-1, checkDigitLength)
	return isbn
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
