package incorrectattr

import "fmt"

type Error struct {
	item string
}

func (i *Error) Error() string {
	return fmt.Sprintf("item: %s", i.item)
}

func New(item string) *Error {
	return &Error{item: item}
}
