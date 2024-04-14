package model

type Document struct {
	Title              string   `json:"title"`
	DocNum             string   `json:"docno"`
	Segments           []string `json:"seg"`
	IdsForIdenticalDoc []string `json:"ids"`
}

func (d *Document) SegLength() int {
	return len(d.Segments)
}

type DocWithSim struct {
	Similarity float64
	DocId      int
}

type DocSimWrapper struct {
	DocId int    // 当前doc的id，从0开始，依次递增
	Title string // 当前doc的title
	Docs  []*DocWithSim
}
