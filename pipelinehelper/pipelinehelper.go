package pipelinehelper

import (
	"bufio"
	"fmt"
	"log"
	"module26/ringbuffer"
	"os"
	"strconv"
	"strings"
	"time"
)

func DataSource(nextStage chan<- int, done chan bool) {
	log.Println("Start getting data")
	scanner := bufio.NewScanner(os.Stdin)
	var data string
	for scanner.Scan() {
		data = scanner.Text()
		if strings.EqualFold(data, "exit") {
			log.Println("Finished working")
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
	log.Println("Start Negative Filter Stage")
	for {
		select {
		case data := <-previousStageChannel:
			log.Println("Inside Negative Filter Stage")
			if data > 0 {
				nextStageChannel <- data
			}
		case <-done:
			return
		}
	}
}

func SpecialFilterStageInt(previousStageChannel <-chan int, nextStageChannel chan<- int, done chan bool) {
	log.Println("Start Special Filter Stage")
	for {
		select {
		case data := <-previousStageChannel:
			log.Println("Inside Special Filter Stage")
			if data%3 == 0 {
				nextStageChannel <- data
			}
		case <-done:
			return
		}
	}
}

func BufferStageInt(previousStageChannel <-chan int, nextStageChannel chan<- int, done chan bool, size int, interval time.Duration) {
	log.Println("Start Buffer Stage")
	buffer := ringbuffer.NewRingBuffer(size)
	for {
		select {
		case data := <-previousStageChannel:
			buffer.Push(data)
		case <-time.After(interval):
			bufferData := buffer.Get()
			if bufferData != nil {
				for _, data := range bufferData {
					nextStageChannel <- data
				}
			}
		case <-done:
			return
		}
	}
}
