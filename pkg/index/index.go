package index

import (
	"godsend/pkg/crawler"
)

type Interface interface {
	Add([]crawler.Document)
	Search(string) []int
	Docs() []crawler.Document
}
