package bestmatch25plus

import (
	"fmt"
	"graph-clustering/model"
	"sort"
)
import "math"

const (
	MAX_CACHE_SIZE = 100
)

// BestMatch25Plus 考虑用泛型替代float32
// 另外如何建立一个随着时间衰减，或者说拥有过时机制的bm25plus 机器呢？这里考虑在流式场景下的检索
type BestMatch25Plus struct {
	cache chan *model.Document
	done  chan struct{}

	k1    float64
	b     float64
	delta float64

	docNum int
	docs   []*model.Document

	idf        map[string]float64 // 每个词的idf值
	df         map[string]int     // 每个词出现该词的文档数量
	wf         []map[string]int   // 文档i里某个词出现的次数
	term2docId map[string][]int   // 某个词对应的包含该词的文档id列表
	totalTerms int                // 所有的文档包含的term的总数
	avgTermCnt float64
}

func NewBm25PlusEngine() *BestMatch25Plus {
	return &BestMatch25Plus{
		k1:    1.2,
		b:     0.75,
		delta: 1.0,

		cache:      make(chan *model.Document, MAX_CACHE_SIZE),
		done:       make(chan struct{}),
		docs:       make([]*model.Document, 0),
		idf:        make(map[string]float64),
		df:         make(map[string]int),
		wf:         make([]map[string]int, 0),
		term2docId: make(map[string][]int),
	}
}

func (bm *BestMatch25Plus) Build() {
	for {
		select {
		case doc := <-bm.cache:
			bm.Receive(doc)
		case <-bm.done:
			bm.avgTermCntCalc()
			bm.IdfCalc()
			fmt.Println("Build BM25+ Engine successful")
			return
		}
	}
}

func (bm *BestMatch25Plus) Receive(doc *model.Document) {
	//fmt.Printf("receive: %v\n", doc)
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

// IdfCalc idf = log(N + 1) - log(n + 0.5), here N + 1 to avoid negative value
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
		bm.term2docId[term] = append(docIdList, bm.docNum)
		return
	}
	bm.term2docId[term] = []int{bm.docNum}
}

func (bm *BestMatch25Plus) ScoreCalc(doc *model.Document, idx int) float64 {
	score := 0.0
	for _, word := range doc.Segments {
		if _, ok := bm.wf[idx][word]; !ok {
			continue
		}
		d := float64(bm.docs[idx].SegLength())
		iwf := float64(bm.wf[idx][word])
		molecular := iwf * (bm.k1 + 1)
		denominator := iwf + bm.k1*(1-bm.b+bm.b*d/bm.avgTermCnt)
		score += bm.idf[word] * (molecular/denominator + bm.delta)
	}
	return score
}

func (bm *BestMatch25Plus) Search(doc *model.Document, topK int) []*model.DocWithSim {
	searched := make(map[int]bool)
	res := make([]*model.DocWithSim, 0)
	for _, term := range doc.Segments {
		for _, docId := range bm.term2docId[term] {
			if !searched[docId] {
				score := bm.ScoreCalc(doc, docId)
				curr := &model.DocWithSim{
					Similarity: score,
					DocId:      docId,
				}
				res = append(res, curr)
				searched[docId] = true
			}
		}
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].Similarity > res[j].Similarity {
			return true
		}
		return false
	})

	res2 := NormalizeAndTruncate(res)

	if topK >= len(res2) {
		return res2
	}
	return res2[:topK]
}

func NormalizeAndTruncate(docs []*model.DocWithSim) []*model.DocWithSim {
	threshold := 0.8
	if len(docs) < 1 {
		return docs
	}
	dominator := docs[0].Similarity
	res := make([]*model.DocWithSim, 0)
	for i := 0; i < len(docs); i++ {
		ratio := docs[i].Similarity / dominator
		if ratio >= threshold {
			res = append(res, docs[i])
		} else {
			break
		}
	}
	return res
}

func (bm *BestMatch25Plus) GetCacheChannel() chan *model.Document {
	return bm.cache
}

func (bm *BestMatch25Plus) GetDoneChannel() chan struct{} {
	return bm.done
}

func (bm *BestMatch25Plus) GetDocByDocId(docId int) *model.Document {
	return bm.docs[docId]
}

func (bm *BestMatch25Plus) GetDocNum() int {
	return bm.docNum
}
