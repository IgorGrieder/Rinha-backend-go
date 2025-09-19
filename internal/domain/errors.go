package domain

import "fmt"

type JSONParsingError struct {
	message string
}

func (e *JSONParsingError) Error() string {
	return fmt.Sprintf("Error while parsing the json of the following string: %s", e.message)
}
