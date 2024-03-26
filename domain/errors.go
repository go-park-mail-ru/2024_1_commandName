package domain

import "fmt"

type CustomError struct {
	Type    string
	Message string
	Segment string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s error in %s: %s", e.Segment, e.Message)
}
