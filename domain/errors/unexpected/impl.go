package unexpected

type Error struct {
	cause error
}

func (i *Error) Error() string {
	return i.cause.Error()
}

func New(cause error) *Error {
	return &Error{cause: cause}
}
