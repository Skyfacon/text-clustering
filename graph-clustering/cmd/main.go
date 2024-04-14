package main

import (
	"fmt"
	"graph-clustering/batchprocess"
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
	fileAddr := "../data/sougouCA_clean_simple.txt"
	dataFetcher := utils.NewDataHandler(fileAddr)
	bm25plus := bestmatch25plus.NewBm25PlusEngine()
	go dataFetcher.Start(bm25plus.GetCacheChannel(), bm25plus.GetDoneChannel())
	bm25plus.Build()

	fmt.Println("Starting BatchProcessing")
	// TODO Search的时候，为避免过短的搜索项，考虑只搜索有用分词个数 >= 3个的
	fmt.Printf("Totally %d docs will be processed\n", len(bm25plus.DocsForSearch))
	batchProcessor := batchprocess.NewBatchProcessor[*model.Document, *model.DocSimWrapper](bm25plus.DocsForSearch, 30)
	batchProcessor.Process(bm25plus.Search)
	res := batchProcessor.GetProcessResult()

	fmt.Println("Start Sort and Show")
	sort.Slice(res, func(i, j int) bool {
		if len(res[i].Docs) > len(res[j].Docs) {
			return true
		}
		return false
	})

	for _, docSearchRes := range res[:100] {
		bm25plus.Show(docSearchRes)
	}

}
