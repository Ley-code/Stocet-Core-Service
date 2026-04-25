package services

// FundamentalsValidationError is returned when input fails business rules
// (negative assets, invalid dates, etc.). HTTP handlers map this to 400.
type FundamentalsValidationError struct {
	Message string
}

func (e *FundamentalsValidationError) Error() string {
	return e.Message
}

// NewFundamentalsValidationError wraps a validation message for fundamentals domain rules.
func NewFundamentalsValidationError(msg string) error {
	return &FundamentalsValidationError{Message: msg}
}
