package isbn

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	versionXParts  = 1
	version10Parts = 4
	version13Parts = 5

	Version10 = 10
	Version13 = 13
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
	isbnRegex = regexp.MustCompile(`(\d+|X|x)`)
)

var (
	errWrongISBN = errors.New("wrong input ISBN format")
)

var (
	version13Weights = []int{1, 3}
)

type ISBN struct {
	version           int
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

func (isbn ISBN) Check() bool {
	if isbn.err != nil || len(isbn.checkDigit) != 1 {
		return false
	}

	// version 13
	if isbn.version == Version13 {
		return strconv.Itoa(isbn.calculateV13CheckDigit()) == isbn.checkDigit
	}

	// version 10 till the end
	v := isbn.calculateV10CheckDigit()
	if v < 10 {
		return strconv.Itoa(v) == isbn.checkDigit
	}

	return "X" == isbn.checkDigit || "x" == isbn.checkDigit
}

func (isbn ISBN) Normalize() {
	if isbn.err != nil || isbn.version == Version13 {
		return
	}

}

func (isbn ISBN) Error() error {
	return isbn.err
}

func (isbn ISBN) String() string {
	// version 13
	if isbn.version == Version13 {
		return fmt.Sprintf("ISBN %s-%s-%s-%s-%s",
			isbn.prefix,
			isbn.registrationGroup,
			isbn.registrant,
			isbn.publication,
			isbn.checkDigit)
	}

	// version 10
	return fmt.Sprintf("ISBN %s-%s-%s-%s",
		isbn.registrationGroup,
		isbn.registrant,
		isbn.publication,
		isbn.checkDigit)
}

// ------------------------------------------------- PRIVATE METHODS -------------------------------------------------

func (isbn ISBN) calculateV13CheckDigit() int {
	w := weightV13()
	sum := weightSum(isbn.prefix, w)
	sum += weightSum(isbn.registrationGroup, w)
	sum += weightSum(isbn.registrant, w)
	sum += weightSum(isbn.publication, w)

	reminder := sum % 10
	if reminder == 0 {
		reminder = 10
	}

	return 10 - reminder
}

func (isbn ISBN) calculateV10CheckDigit() int {
	w := weightV10()
	sum := weightSum(isbn.registrationGroup, w)
	sum += weightSum(isbn.registrant, w)
	sum += weightSum(isbn.publication, w)

	// reminder
	return 11 - (sum % 11)
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
		// TODO - there is only one big number - wihtout hyphnes/spaces
		fallthrough // remove when implemented
	default:
		return ISBN{err: errWrongISBN}
	}
}

func weightV13() func() int {
	idx := -1
	return func() int {
		idx++
		return version13Weights[idx%2]

	}
}

func weightV10() func() int {
	value := 10
	return func() int {
		v := value
		value--
		return v
	}
}

func weightSum(number string, weight func() int) int {
	sum := 0
	for _, v := range number {
		sum += int(v-'0') * weight()
	}

	return sum
}
