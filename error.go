package trapi

import "fmt"

type ParserError struct {
	Message  string
	Filename string
	Line     int
}

func (e *ParserError) Error() string {
	if e.Filename == "" {
		return e.Message
	}
	return fmt.Sprintf("%s [%s:%d]", e.Message, e.Filename, e.Line)
}

func NewParserError(message string, filename string, line int) *ParserError {
	return &ParserError{
		Message:  message,
		Filename: filename,
		Line:     line,
	}
}
