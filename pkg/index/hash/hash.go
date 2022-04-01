package hash

import (
	"godsend/pkg/crawler"
	"sort"
	"strings"
)

type Index struct {
	lexems map[string][]int
	docs   []crawler.Document
}

func New() *Index {
	var i Index
	i.lexems = make(map[string][]int)
	i.docs = make([]crawler.Document, 0)
	return &i
}

func (i *Index) Add(docs []crawler.Document) {

	// sort by ID
	sort.Slice(docs[:], func(i, j int) bool {
		return docs[i].ID < docs[j].ID
	})

	for _, doc := range docs {
		// add docs to index
		i.docs = append(i.docs, doc)

		// index lexems
		for _, word := range strings.Split(doc.Title, " ") {
			w := strings.ToLower(word)

			if !present(i.lexems[w], doc.ID) {
				i.lexems[w] = append(i.lexems[w], doc.ID)
			}
		}
	}
}

func (i *Index) Search(word string) []int {
	return i.lexems[strings.ToLower(word)]
}

func (i *Index) Docs() []crawler.Document {
	return i.docs
}

func present(ids []int, itemId int) bool {
	for _, id := range ids {
		if id == itemId {
			return true
		}
	}
	return false
}
