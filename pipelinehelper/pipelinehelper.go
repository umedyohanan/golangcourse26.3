package pipelinehelper

import (
	"bufio"
	"fmt"
	"module20/ringbuffer"
	"os"
	"strconv"
	"strings"
	"time"
)

func DataSource(nextStage chan<- int, done chan bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var data string
	for scanner.Scan() {
		data = scanner.Text()
		if strings.EqualFold(data, "exit") {
			fmt.Println("Finished working")
			close(done)
			return
		}
		i, err := strconv.Atoi(data)
		if err != nil {
			fmt.Println("Only digits are available")
			continue
		}
		nextStage <- i
	}
}

func NegativeFilterStageInt(previousStageChannel <-chan int, nextStageChannel chan<- int, done chan bool) {
	for {
		select {
		case data := <- previousStageChannel:
			if data > 0 {
				nextStageChannel <- data
			}
		case <- done:
			return
		}
	}
}

func SpecialFilterStageInt(previousStageChannel <-chan int, nextStageChannel chan<- int, done chan bool) {
	for {
		select {
		case data := <-previousStageChannel:
			if data%3 == 0 {
				nextStageChannel <- data
			}
		case <- done:
			return
		}
	}
}

func BufferStageInt(previousStageChannel <-chan int, nextStageChannel chan<- int, done chan bool, size int, interval time.Duration) {
	buffer := ringbuffer.NewRingBuffer(size)
	for {
		select {
		case data := <-previousStageChannel:
			buffer.Push(data)
		case <- time.After(interval):
			bufferData := buffer.Get()
			if bufferData != nil {
				for _, data := range bufferData {
					nextStageChannel <- data
				}
			}
		case <- done:
			return
		}
	}
}