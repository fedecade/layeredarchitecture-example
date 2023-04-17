package translator

type impl struct{}

func New() Translator {
	return &impl{}
}
