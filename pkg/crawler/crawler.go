package crawler

type Interface interface {
	Scan(url string, depth int) ([]Document, error)
}

type Document struct {
	ID    int
	URL   string
	Title string
}
