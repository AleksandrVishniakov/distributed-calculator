package signs

import "slices"

var (
	operations      = []byte("+-*/")
	digits          = []byte("0123456789")
	brackets        = []byte("()")
	space      byte = ' '
)

func IsOperation(b byte) bool {
	return slices.Contains(operations, b)
}

func IsDigit(b byte) bool {
	return slices.Contains(digits, b)
}

func IsBracket(b byte) bool {
	return slices.Contains(brackets, b)
}

func IsSpace(b byte) bool {
	return b == space
}

func IsValidSymbol(b byte) bool {
	return IsDigit(b) || IsBracket(b) || IsSpace(b) || IsOperation(b)
}
