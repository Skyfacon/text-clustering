package model

type Document struct {
	Title    string   `json:"title"`
	DocNum   string   `json:"docno"`
	Segments []string `json:"seg"`
}

func (d *Document) SegLength() int {
	return len(d.Segments)
}

type DocWithSim struct {
	Similarity float64
	DocId      int
}
