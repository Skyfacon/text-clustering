package batchprocess

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

type batchProcess[T interface{}, S interface{}] struct {
	items     []T
	batchSize int
	batchCnt  int
	result    []S
}

func NewBatchProcessor[T interface{}, S interface{}](items []T, batchCnt int) *batchProcess[T, S] {
	batchSize := divideAndFloor(len(items), batchCnt)
	fmt.Printf("batchSize:%d\n", batchSize)
	return &batchProcess[T, S]{
		items:     items,
		batchSize: batchSize,
		batchCnt:  batchCnt,
		result:    make([]S, len(items)),
	}

}

func (p *batchProcess[T, S]) getBatch(batch int) ([]T, error) {
	if batch > p.batchCnt || batch <= 0 {
		return nil, errors.New(fmt.Sprintf("batch is invalid, should be 0 < batch <= %d", p.batchSize))
	}
	if batch == p.batchCnt {
		return p.items[(batch-1)*p.batchSize:], nil
	}
	return p.items[(batch-1)*p.batchSize : min(len(p.items), batch*p.batchSize)], nil
}

func (p *batchProcess[T, S]) Process(processSingle func(item T) S) {
	var wg sync.WaitGroup
	for i := 1; i <= p.batchCnt; i++ {
		currBatch, _ := p.getBatch(i)
		wg.Add(1)
		go func(batchItem []T, batch int) {
			for cidx, item := range batchItem {
				currRes := processSingle(item)
				p.result[(batch-1)*p.batchSize+cidx] = currRes
			}
			wg.Done()
		}(currBatch, i)
	}
	wg.Wait()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func divideAndFloor(a, b int) int {
	quotient := float64(a) / float64(b)
	return int(math.Floor(quotient))
}
