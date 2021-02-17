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
	diff := time.Now().Sub(startTime).Milliseconds()
	fmt.Printf("load file:%s, count:%d, cost:%dms\n", path, len(docs), diff)

	matchList := []matchTerm{
		matchTermContain,
		matchTermRegexp,
	}
	for _, m := range matchList {
		hitDocs := search(docs, *term, m)
		printResult(hitDocs)
	}
}

type Document struct {
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Text  string `xml:"abstract"`
	ID    int
}

func loadDucments(path string) ([]Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	dec := xml.NewDecoder(f)
	dump := struct {
		Documents []Document `xml:"doc"`
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

func search(docs []Document, term string, m matchTerm) []Document {
	startTime := time.Now()
	defer func() {
		ms := time.Now().Sub(startTime).Milliseconds()
		funName := runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
		fmt.Printf("search, func:%s, cost:%dms\n", funName, ms)
	}()

	r := make([]Document, 0, len(docs))
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

func printResult(docs []Document) {
	fmt.Printf("hit results count:%d\n", len(docs))

	for _, doc := range docs {
		fmt.Println(doc)
	}
}
