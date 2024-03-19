package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"graph-clustering/model"
	"os"
	"time"
)

const (
	MAX_DATA_BUFFER_SIZE = 100
)

// RawData 无论是哪种格式的数据，都需要format成该标准格式后，再喂给聚类算法
type RawData struct {
	Title   string   `json:"title"`
	DocNum  string   `json:"docno"`
	SegList []string `json:"seg"`
}

type DataHandler struct {
	filePath string
	drainer  chan *RawData
	err      chan error
	done     chan struct{}
}

type Runner interface {
	Start() error
	Stop() error
}

func NewDataHandler(filePath string) *DataHandler {
	return &DataHandler{
		filePath: filePath,
		drainer:  make(chan *RawData, MAX_DATA_BUFFER_SIZE),
		err:      make(chan error),
		done:     make(chan struct{}),
	}
}

func (h *DataHandler) read(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		h.err <- errors.New("open file error: " + err.Error())
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	// 创建一个Scanner，用于逐行读取文件
	scanner := bufio.NewScanner(file)

	// 循环读取文件中的每一行
	for scanner.Scan() {
		// 获取当前行的文本内容
		line := scanner.Text()

		// 解析JSON字符串到结构体
		var data RawData
		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			fmt.Printf("Error parsing JSON: %v", err)
			continue
		}
		h.drainer <- &data
	}

	cnt := 0
	for {
		if len(h.drainer) == 0 {
			h.done <- struct{}{}
			return
		}
		time.Sleep(1 * time.Second)
		cnt++
		fmt.Printf("wait %ds\n", cnt)
	}

}

// Start 可以传入一个function，来处理h.drainer通道的内容
func (h *DataHandler) Start(cache chan *model.Document, done chan struct{}) error {
	go h.read(h.filePath)
	for {
		select {
		case data := <-h.drainer:
			cache <- &model.Document{
				Title:    data.Title,
				DocNum:   data.DocNum,
				Segments: data.SegList,
			}
		case err := <-h.err:
			fmt.Printf("%v\n", err)
			h.Stop()
			return err
		case <-h.done:
			done <- struct{}{}
			h.Stop()
			return nil
		}
	}
}

func (h *DataHandler) Stop() error {
	close(h.drainer)
	close(h.err)
	close(h.done)
	return nil
}
