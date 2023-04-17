package customer

type Repository interface {
	Register(Customer) error
}
