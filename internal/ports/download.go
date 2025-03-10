package ports

type Downloader interface {
	Download(url string) (string, error)
}
