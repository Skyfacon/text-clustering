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

func (bm *BestMatch25Plus) Show(docs []*model.DocWithSim) {
	for _, doc := range docs {
		title := bm.docs[doc.DocId].Title
		score := doc.Similarity
		segs := bm.docs[doc.DocId].Segments
		fmt.Printf("%s, %.2f,%v\n", title, score, segs)
	}
}

func (bm *BestMatch25Plus) ShowStructure() {
	fmt.Println("docs:")
	for _, doc := range bm.docs {
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
