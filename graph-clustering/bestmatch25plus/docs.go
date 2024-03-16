package bestmatch25plus

type Document struct {
	Title    string   `json:"title"`
	DocNum   string   `json:"docno"`
	Segments []string `json:"seg"`
}

func (d *Document) SegLength() int {
	return len(d.Segments)
}
