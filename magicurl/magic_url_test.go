package magicurl

import (
	"fmt"
	"testing"
)

func TestBase62Encoding(t *testing.T) {
	assertCorrectBase62Encoding := func(t *testing.T, base10, expectedBase62 string) {
		t.Helper()
		actual, _ := EncodeToBase62(base10)
		if expectedBase62 != actual {
			t.Errorf("Expected:%s does not match the actual: %s\n", expectedBase62, actual)
		}
	}

	t.Run("single digit decimal", func(t *testing.T) {
		assertCorrectBase62Encoding(t, "0", "A")
	})
	t.Run("digit decimal", func(t *testing.T) {
		assertCorrectBase62Encoding(t, "62", "BA")
	})
	t.Run("five digit decimal", func(t *testing.T) {
		assertCorrectBase62Encoding(t, "85938", "WWG")
	})
	t.Run("seven digit decimal", func(t *testing.T) {
		assertCorrectBase62Encoding(t, "8789312", "k2fG")
	})
	t.Run("nine digit decimal", func(t *testing.T) {
		assertCorrectBase62Encoding(t, "2500000000", "CtLuMo")
	})
	t.Run("thirteen digit decimal", func(t *testing.T) {
		assertCorrectBase62Encoding(t, "3521614600000", "99998X2")
	})
	t.Run("max size uint64", func(t *testing.T) {
		assertCorrectBase62Encoding(t, "18446744073709551615", "V8qRkBGKRiP")
	})
	t.Run("invalid decimal input returns magicUrlBaseConverterError", func(t *testing.T) {
		if _, err := EncodeToBase62("100+s0000"); err != nil {
			fmt.Println("TEST CAUGHT ERR")
			if e, caught := err.(*magicURLBaseConverterError); !caught {
				t.Errorf("Expected magicUrlBaseConverterError, but got: %s\n", e)
			}
		}
	})
}

func TestDecodeToBase10(t *testing.T) {
	assertCorrectBase10Decoding := func(t *testing.T, base10, expectedBase62 string) {
		t.Helper()
		actual, _ := DecodeToBase10(base10)
		if expectedBase62 != actual {
			t.Errorf("Expected:%s does not match the actual: %s\n", expectedBase62, actual)
		}
	}
	t.Run("encode base62 single digit decimal", func(t *testing.T) {
		assertCorrectBase10Decoding(t, "A", "0")
	})
	t.Run("encode base62 two digit decimal", func(t *testing.T) {
		assertCorrectBase10Decoding(t, "BA", "62")
	})
	t.Run("encode base62 five digit decimal", func(t *testing.T) {
		assertCorrectBase10Decoding(t, "WWG", "85938")
	})
	t.Run("encode base62 seven digit decimal", func(t *testing.T) {
		assertCorrectBase10Decoding(t, "k2fG", "8789312")
	})
	t.Run("encode base62 with nine digit decimal", func(t *testing.T) {
		assertCorrectBase10Decoding(t, "CtLuMo", "2500000000")
	})
	t.Run("encode base62 with thirteen digit decimal", func(t *testing.T) {
		assertCorrectBase10Decoding(t, "99998X2", "3521614600000")
	})
	t.Run("encode base62 with max size uint64", func(t *testing.T) {
		assertCorrectBase10Decoding(t, "V8qRkBGKRiP", "18446744073709551615")
	})
}
