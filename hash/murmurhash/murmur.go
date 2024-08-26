package murmurhash

import (
	"fmt"
	"strings"
)

const (
	base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func decimalToBase62(decimal int) string {
	if decimal == 0 {
		return string(base62Alphabet[0])
	}

	base := len(base62Alphabet)
	result := ""
	for decimal > 0 {
		remainder := decimal % base
		result = string(base62Alphabet[remainder]) + result
		decimal /= base
	}

	return result
}

func base62ToDecimal(base62 string) int {
	base := len(base62Alphabet)
	result := 0
	for _, char := range base62 {
		result = result*base + strings.IndexRune(base62Alphabet, char)
	}

	return result
}

func main() {
	decimalNumber := 12345
	base62Number := decimalToBase62(decimalNumber)
	fmt.Printf("Decimal %d in Base62: %s\n", decimalNumber, base62Number)

	convertedBack := base62ToDecimal(base62Number)
	fmt.Printf("Base62 %s in Decimal: %d\n", base62Number, convertedBack)
}
