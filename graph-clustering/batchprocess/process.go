package batchprocess

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

// BatchProcess TODO: 这里没有充分利用好CPU，运行时长完全取决于最慢的那个Batch，考虑采用进程池，小batch的方式，来进行加速
type BatchProcess[T interface{}, S interface{}] struct {
	items     []T
	batchSize int
	batchCnt  int
	result    []S
}

func NewBatchProcessor[T interface{}, S interface{}](items []T, batchCnt int) *BatchProcess[T, S] {
	batchSize := divideAndFloor(len(items), batchCnt)
	fmt.Printf("batchSize:%d\n", batchSize)
	return &BatchProcess[T, S]{
		items:     items,
		batchSize: batchSize,
		batchCnt:  batchCnt,
		result:    make([]S, len(items)),
	}

}

func (p *BatchProcess[T, S]) getBatch(batch int) ([]T, error) {
	if batch > p.batchCnt || batch <= 0 {
		return nil, errors.New(fmt.Sprintf("batch is invalid, should be 0 < batch <= %d", p.batchSize))
	}
	if batch == p.batchCnt {
		return p.items[(batch-1)*p.batchSize:], nil
	}
	return p.items[(batch-1)*p.batchSize : min(len(p.items), batch*p.batchSize)], nil
}

func (p *BatchProcess[T, S]) Process(processSingle func(item T) S) {
	var wg sync.WaitGroup
	for i := 1; i <= p.batchCnt; i++ {
		currBatch, _ := p.getBatch(i)
		wg.Add(1)
		go func(batchItem []T, batch int) {
			defer wg.Done()
			for idx, item := range batchItem {
				currRes := processSingle(item)
				p.result[(batch-1)*p.batchSize+idx] = currRes
			}
		}(currBatch, i)
	}
	wg.Wait()
}

func (p *BatchProcess[T, S]) GetProcessResult() []S {
	return p.result
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
