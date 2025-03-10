package ports

type Persister interface {
	Persist(content string) (string, error)
}
