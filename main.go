package main

import (
	"fmt"
	"log"
	"module26/pipelinehelper"
	"os"
	"os/signal"
	"time"
)

const bufferDrainInterval time.Duration = 10 * time.Second

const bufferSize int = 5

func main() {
	input := make(chan int)
	done := make(chan bool)
	go pipelinehelper.DataSource(input, done)

	negativeFilterChan := make(chan int)
	go pipelinehelper.NegativeFilterStageInt(input, negativeFilterChan, done)

	specialFilterChan := make(chan int)
	go pipelinehelper.SpecialFilterStageInt(negativeFilterChan, specialFilterChan, done)

	bufferedIntChannel := make(chan int)
	go pipelinehelper.BufferStageInt(specialFilterChan, bufferedIntChannel, done, bufferSize, bufferDrainInterval)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	for {
		select {
		case <-sig:
			log.Println("Got os signal, stopping")
			return
		case data := <-bufferedIntChannel:
			fmt.Println("Produced data, ", data)
		case <-done:
			return
		}
	}
}
