package domain

import "fmt"

type CustomError struct {
	Type    string
	Message string
	Segment string
}

func (e CustomError) Error() string {
	//TODO (fix incorrect) res := fmt.Sprintf("error in %s: %s", e.Segment, e.Message)
	res := fmt.Sprintf("%s", e.Message)
	return res
}
