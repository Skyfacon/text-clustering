package main

import (
	"fmt"
	"graph-clustering/bestmatch25plus"
	"graph-clustering/model"
	"graph-clustering/utils"
	"sort"
)

type R struct {
	title  string
	docs   []*model.DocWithSim
	docNum int
}

func main() {
	fileAddr := "../data/sougouCA_seg_simple.txt"
	dataFetcher := utils.NewDataHandler(fileAddr)
	bm25plus := bestmatch25plus.NewBm25PlusEngine()
	go dataFetcher.Start(bm25plus.GetCacheChannel(), bm25plus.GetDoneChannel())
	bm25plus.Build()

	res := make([]*R, 0)
	for i := 0; i < bm25plus.GetDocNum(); i++ {
		doc := bm25plus.GetDocByDocId(i)
		r := bm25plus.Search(doc, 100)

		res = append(res, &R{
			title:  doc.Title,
			docs:   r,
			docNum: len(r),
		})
	}

	sort.Slice(res, func(i, j int) bool {
		if res[i].docNum > res[j].docNum {
			return true
		}
		return false
	})

	for _, doc := range res[:100] {
		fmt.Printf(doc.title + "\n")
		bm25plus.Show(doc.docs)
		fmt.Printf("-------------------------\n")
	}

}
