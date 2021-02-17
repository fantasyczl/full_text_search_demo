package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/fantasyczl/full_text_search_demo/fts"
)

func main() {
	fmt.Println("app starting...")

	term := flag.String("term", "", "term for search")
	flag.Parse()

	if *term == "" {
		fmt.Println("no term argument")
		os.Exit(2)
	}

	startTime := time.Now()
	path := "./data/dump.xml"
	docs, err := loadDucments(path)
	if err != nil {
		fmt.Println("load documents failed, ", err)
		os.Exit(2)
	}
	diff := time.Now().Sub(startTime).Microseconds()
	fmt.Printf("load file:%s, count:%d, cost:%dus\n", path, len(docs), diff)

	// search
	matchList := []matchTerm{
		matchTermContain,
		matchTermRegexp,
	}
	for _, m := range matchList {
		hitDocs := search(docs, *term, m)
		printResult(hitDocs)
	}

	// build index
	idx := fts.BuildIndex(docs)
	startTime = time.Now()
	docIds := idx.Search(*term)
	diff = time.Now().Sub(startTime).Microseconds()
	fmt.Println(docIds, "cost:", diff, "us")
}

func loadDucments(path string) ([]fts.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	dec := xml.NewDecoder(f)
	dump := struct {
		Documents []fts.Document `xml:"doc"`
	}{}

	if err := dec.Decode(&dump); err != nil {
		return nil, err
	}

	docs := dump.Documents
	for i := range docs {
		docs[i].ID = i
	}

	return docs, nil
}

func search(docs []fts.Document, term string, m matchTerm) []fts.Document {
	startTime := time.Now()
	defer func() {
		ms := time.Now().Sub(startTime).Microseconds()
		funName := runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
		fmt.Printf("search, func:%s, cost:%dus\n", funName, ms)
	}()

	r := make([]fts.Document, 0, len(docs))
	for _, doc := range docs {
		if m(doc.Text, term) {
			r = append(r, doc)
		}
	}

	return r
}

type matchTerm func(string, string) bool

func matchTermContain(s string, term string) bool {
	return strings.Contains(s, term)
}

func matchTermRegexp(s string, term string) bool {
	re := regexp.MustCompile(fmt.Sprintf(`(?i)\b%s\b`, term))
	return re.MatchString(s)
}

func printResult(docs []fts.Document) {
	fmt.Printf("hit results count:%d\n", len(docs))

	for _, doc := range docs {
		fmt.Println(doc)
	}
}
