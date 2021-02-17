package fts

import (
	"fmt"
	"time"
)

type Index map[string][]int

func (idx Index) Add(docs []Document) {
	for _, doc := range docs {
		for _, token := range Analyze(doc.Text) {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == doc.ID {
				continue
			}

			idx[token] = append(ids, doc.ID)
		}
	}
}

func (idx Index) Search(text string) []int {
	var r []int
	for _, token := range Analyze(text) {
		if ids, ok := idx[token]; ok {
			if r == nil {
				r = ids
			} else {
				r = Intersection(r, ids)
			}
		} else {
			// token doesn't exist
			return nil
		}
	}

	return r
}

func BuildIndex(docs []Document) Index {
	startTime := time.Now()
	defer func() {
		diff := time.Now().Sub(startTime).Microseconds()
		fmt.Printf("buildIndex cost:%dus\n", diff)
	}()

	idx := make(Index)
	idx.Add(docs)

	return idx
}

func Intersection(a []int, b []int) []int {
	maxLen := len(a)
	if len(b) > len(a) {
		maxLen = len(b)
	}

	r := make([]int, 0, maxLen)

	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			r = append(r, a[i])
			i++
			j++
		}
	}

	return r
}
