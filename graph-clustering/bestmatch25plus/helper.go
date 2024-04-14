package bestmatch25plus

import (
	"fmt"
	"github.com/yanyiwu/gojieba"
	"graph-clustering/model"
	"strings"
	"time"
)

func dictFunc(words []string) {
	dict := make(map[string]int)
	for _, word := range words {
		dict[word] += 1
	}
	for k, v := range dict {
		fmt.Printf("k:%s, v:%d\n", k, v)
	}
}

func (bm *BestMatch25Plus) ShowDocSim(docs []*model.DocWithSim) {
	for _, doc := range docs {
		title := bm.Docs[doc.DocId].Title
		score := doc.Similarity
		segs := bm.Docs[doc.DocId].Segments
		idCnt := len(bm.Docs[doc.DocId].IdsForIdenticalDoc) + 1
		fmt.Printf("%s,%.2f,%d,%v\n", title, score, idCnt, segs)
	}
}

func (bm *BestMatch25Plus) Show(docSearchResult *model.DocSimWrapper) {
	fmt.Printf(docSearchResult.Title + "\n")
	for _, doc := range docSearchResult.Docs {
		title := bm.Docs[doc.DocId].Title
		score := doc.Similarity
		idCnt := len(bm.Docs[doc.DocId].IdsForIdenticalDoc) + 1
		fmt.Printf("%s,%.2f,%d\n", title, score, idCnt)
	}
}

func (bm *BestMatch25Plus) ShowStructure() {
	fmt.Println("Docs:")
	for _, doc := range bm.Docs {
		fmt.Printf(doc.Title + "---" + strings.Join(doc.Segments, "|") + "\n")
	}

	fmt.Println("wf:")
	for word, freq := range bm.df {
		fmt.Printf("%s:%d\n", word, freq)
	}

	fmt.Println("word2docId:")
	for word, arr := range bm.term2docId {
		fmt.Printf("%s:%v\n", word, arr)
	}

	fmt.Printf("docNum:%d,totalTerms:%d,avgTermCnt:%f\n", bm.docNum, bm.totalTerms, bm.avgTermCnt)
}

func simulate(sentences []string) {
	x := gojieba.NewJieba()
	defer x.Free()

	bm25 := NewBm25PlusEngine()
	go bm25.Build()

	for i, s := range sentences {
		curr := x.Cut(s, true)
		doc := &model.Document{
			Title:    s,
			DocNum:   fmt.Sprintf("%d", i),
			Segments: curr,
		}
		bm25.cache <- doc
		time.Sleep(time.Second)
	}
	bm25.done <- struct{}{}

	bm25.ShowStructure()
}

func wordSegInGo(s string) {
	x := gojieba.NewJieba()
	defer x.Free()

	words := x.Cut(s, true)
	fmt.Printf(strings.Join(words, " ") + "\n")
}
