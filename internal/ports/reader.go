package ports

type Reader interface {
	ReadLines(urlChan chan<- string) error
}
