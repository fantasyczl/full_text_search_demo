package fts

type Document struct {
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Text  string `xml:"abstract"`
	ID    int
}
