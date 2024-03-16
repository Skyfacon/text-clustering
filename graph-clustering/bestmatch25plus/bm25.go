package bestmatch25plus

import "fmt"
import "math"

// BestMatch25Plus 考虑用泛型替代float32
type BestMatch25Plus struct {
	k1    float64
	b     float64
	delta float64

	docNum int
	docs   []*Document

	idf        map[string]float64 // 每个词的idf值
	df         map[string]int     // 每个词出现该词的文档数量
	wf         []map[string]int   // 文档i里某个词出现的次数
	term2docId map[string][]int   // 某个词对应的包含该词的文档id列表
	totalTerms int                // 所有的文档包含的term的总数
	avgTermCnt float64
}

func (bm *BestMatch25Plus) Receive(doc *Document) {
	bm.docs = append(bm.docs, doc)

	wordMap := make(map[string]int)
	for _, word := range doc.Segments {
		wordMap[word] += 1
		// 该词在本文档中第一次出现时
		if wordMap[word] == 1 {
			bm.df[word] += 1
			bm.addTerm(word)
		}
	}

	bm.wf = append(bm.wf, wordMap)
	bm.totalTerms += len(doc.Segments)
	bm.docNum += 1
}

func (bm *BestMatch25Plus) IdfCalc() {
	for k, v := range bm.df {
		bm.idf[k] = math.Log(float64(bm.docNum)+1) - math.Log(float64(v)+0.5)
	}
}

func (bm *BestMatch25Plus) avgTermCntCalc() {
	bm.avgTermCnt = float64(bm.totalTerms) / float64(bm.docNum)
}

func (bm *BestMatch25Plus) addTerm(term string) {
	if docIdList, ok := bm.term2docId[term]; ok {
		docIdList = append(docIdList, bm.docNum)
		return
	}
	bm.term2docId[term] = []int{bm.docNum}
}

func (bm *BestMatch25Plus) ScoreCalc(doc *Document, idx int) float64 {
	score := 0.0
	for _, word := range doc.Segments {
		if _, ok := bm.wf[idx][word]; !ok {
			continue
		}
		d := float64(bm.docs[idx].SegLength())
		iwf := float64(bm.wf[idx][word])
		molecular := bm.delta + iwf*(bm.k1+1)
		denominator := iwf + bm.k1*(1-bm.b+bm.b*d/bm.avgTermCnt)
		score += bm.idf[word] * molecular / denominator
	}
	return score
}

func (bm *BestMatch25Plus) Search() {

}

func dictFunc(words []string) {
	dict := make(map[string]int)
	for _, word := range words {
		dict[word] += 1
	}
	for k, v := range dict {
		fmt.Printf("k:%s, v:%d\n", k, v)
	}
}
