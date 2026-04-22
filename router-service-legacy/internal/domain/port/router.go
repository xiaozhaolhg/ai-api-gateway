package port

type Router interface {
	SelectProvider(model string) (Provider, error)
}