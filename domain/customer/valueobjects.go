package customer

import (
	"strings"

	"example.layeredarch/domain/errors/incorrectattr"
)

type Customer interface {
	Name() string
	Email() string
}

type impl struct {
	name  string
	email string
}

func (i *impl) Name() string {
	return i.name
}

func (i *impl) Email() string {
	return i.email
}

type Builder interface {
	New(
		name string,
		email string,
	) (Customer, error)
}

type builder struct{}

func (b *builder) New(
	name string,
	email string,
) (Customer, error) {
	if len(strings.TrimSpace(name)) == 0 {
		return nil, incorrectattr.New("name")
	}
	if len(strings.TrimSpace(email)) == 0 {
		return nil, incorrectattr.New("email")
	}
	return &impl{
		name:  name,
		email: email,
	}, nil
}

func NewBuilder() Builder {
	return &builder{}
}
