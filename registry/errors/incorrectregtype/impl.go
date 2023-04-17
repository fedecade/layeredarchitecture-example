package incorrectregtype

import (
	"fmt"
)

type Error struct {
	name string
	typ  string
}

func (i *Error) Error() string {
	return fmt.Sprintf(
		"Incorrect Type: %s %s",
		i.name,
		i.typ,
	)
}

func New(name string, typ string) *Error {
	return &Error{name: name, typ: typ}
}
