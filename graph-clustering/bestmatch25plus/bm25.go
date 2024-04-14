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

	docNum        int
	Docs          []*model.Document // 仅包含unique title的doc
	DocsForSearch []*model.Document // 对需要进行检索的 doc添加了一些限制，如分词数量>=3等，保证信息的完整

	idf        map[string]float64 // 每个词的idf值
	df         map[string]int     // 每个词出现该词的文档数量
	wf         []map[string]int   // 文档i里某个词出现的次数
	term2docId map[string][]int   // 某个词对应的包含该词的文档id列表
	totalTerms int                // 所有的文档包含的term的总数
	avgTermCnt float64            // 每篇文档的平均切词数量

	uniqueTitleMap  map[string]int
	docNum2docIdMap map[string]int
}

func NewBm25PlusEngine() *BestMatch25Plus {
	return &BestMatch25Plus{
		k1:    1.2,
		b:     0.75,
		delta: 1.0,

		cache:         make(chan *model.Document, MAX_CACHE_SIZE),
		done:          make(chan struct{}),
		Docs:          make([]*model.Document, 0),
		DocsForSearch: make([]*model.Document, 0),
		idf:           make(map[string]float64),
		df:            make(map[string]int),
		wf:            make([]map[string]int, 0),
		term2docId:    make(map[string][]int),

		uniqueTitleMap:  make(map[string]int),
		docNum2docIdMap: make(map[string]int),
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
	if idx, ok := bm.uniqueTitleMap[doc.Title]; ok {
		bm.Docs[idx].IdsForIdenticalDoc = append(bm.Docs[idx].IdsForIdenticalDoc, doc.DocNum)
		bm.docNum2docIdMap[doc.DocNum] = idx
		return
	}

	bm.Docs = append(bm.Docs, doc)
	bm.uniqueTitleMap[doc.Title] = bm.docNum
	bm.docNum2docIdMap[doc.DocNum] = bm.docNum

	// TODO 需要搜索的文档的过滤逻辑，只做清洗数据之用，后续对于一般文本，可以移除
	if len(doc.Segments) >= 3 {
		bm.DocsForSearch = append(bm.DocsForSearch, doc)
	}

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
		d := float64(bm.Docs[idx].SegLength())
		iwf := float64(bm.wf[idx][word])
		molecular := iwf * (bm.k1 + 1)
		denominator := iwf + bm.k1*(1-bm.b+bm.b*d/bm.avgTermCnt)
		score += bm.idf[word] * (molecular/denominator + bm.delta)
	}
	return score
}

func (bm *BestMatch25Plus) SearchTopK(doc *model.Document, topK int) *model.DocSimWrapper {
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
		return &model.DocSimWrapper{
			DocId: bm.docNum2docIdMap[doc.DocNum],
			Title: doc.Title,
			Docs:  res2,
		}
	}
	return &model.DocSimWrapper{
		DocId: bm.docNum2docIdMap[doc.DocNum],
		Title: doc.Title,
		Docs:  res2[:topK],
	}
}

func (bm *BestMatch25Plus) Search(doc *model.Document) *model.DocSimWrapper {
	topK := 100
	return bm.SearchTopK(doc, topK)
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
	return bm.Docs[docId]
}

func (bm *BestMatch25Plus) GetDocNum() int {
	return bm.docNum
}
