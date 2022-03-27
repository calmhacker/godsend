package crawler

type Interface interface {
	Scan(url string, depth int) ([]Document, error)
	BatchScan(urls []string, depth int, workers int) (<-chan Document, <-chan error)
}

type Document struct {
	ID    int
	URL   string
	Title string
	Body  string
}
