package scraper

type Handler[T any] interface {
	Handle(value T)
}
