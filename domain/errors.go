package domain

import "fmt"

type CustomError struct {
	Type    string
	Message string
	Segment string
}

func (e CustomError) Error() string {
	res := fmt.Sprintf("%s error in %s: %s", e.Type, e.Segment, e.Message)
	return res
}
