package calculator

import "errors"

var (
	ErrDivideByZero = errors.New("calculator: divide by zero")
)

type Calculator interface {
	Add(a, b float64) (float64, error)
	Subtract(a, b float64) (float64, error)
	Divide(a, b float64) (float64, error)
	Multiply(a, b float64) (float64, error)
}

type simpleCalculator struct{}

func NewSimpleCalculator() Calculator {
	return &simpleCalculator{}
}

func (s simpleCalculator) Add(a, b float64) (float64, error) {
	return a + b, nil
}

func (s simpleCalculator) Subtract(a, b float64) (float64, error) {
	return a - b, nil
}

func (s simpleCalculator) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}

	return a / b, nil
}

func (s simpleCalculator) Multiply(a, b float64) (float64, error) {
	return a * b, nil
}
