package magicurl

import (
	"fmt"
	"strings"
)

type magicURLBaseConverterError struct {
	Message string
}

func (m *magicURLBaseConverterError) Error() string {
	return m.Message
}

var base62Characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
var base10Characters = "0123456789"

//EncodeToBase62 takes a base 10 representation and returns the base62 encoding
func EncodeToBase62(numDecimal string) (string, error) {
	err := validateDecimal(numDecimal)
	if err != nil {
		return "", err
	}
	return convert(numDecimal, 10, 62)
}

//DecodeToBase10 takes a base 62 representation and returns the encoding in base 10
func DecodeToBase10(numBase62 string) (string, error) {
	err := validateBase62(numBase62)
	if err != nil {
		return "", err
	}
	return convert(numBase62, 62, 10)
}

func validateDecimal(num string) error {
	return validateNumExistsInCharacterSet(num, base10Characters)
}

func validateBase62(num string) error {
	return validateNumExistsInCharacterSet(num, base62Characters)
}

func validateNumExistsInCharacterSet(num, characterSet string) error {
	for _, ch := range num {
		if !strings.Contains(characterSet, string(ch)) {
			message := fmt.Sprintf("Invalid character %c found for character set: %s", ch, characterSet)
			return &magicURLBaseConverterError{message}
		}
	}
	return nil
}

func convert(num string, fromBase, toBase uint) (string, error) {
	var sourceCharacterSet string
	var targetCharacterSet string

	if fromBase == 10 && toBase == 62 {
		sourceCharacterSet = base10Characters
		targetCharacterSet = base62Characters
	} else if fromBase == 62 && toBase == 10 {
		sourceCharacterSet = base62Characters
		targetCharacterSet = base10Characters
	} else {
		message := fmt.Sprintf("Unsupported base conversion, from base: %d to base: %d", fromBase, toBase)
		return "", &magicURLBaseConverterError{message}
	}

	var decimalRepr uint
	for _, c := range num {
		decimalRepr = decimalRepr*fromBase + uint(strings.IndexRune(sourceCharacterSet, c))
	}
	if decimalRepr == 0 {
		return string(targetCharacterSet[0]), nil
	}

	//TODO: Refactor for string building
	var digits []string
	for decimalRepr != 0 {
		remainder := decimalRepr % toBase
		digits = append(digits, string(targetCharacterSet[remainder]))
		decimalRepr /= toBase
	}

	reverse(digits)
	return strings.Join(digits, ""), nil
}

func reverse(a []string) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}
