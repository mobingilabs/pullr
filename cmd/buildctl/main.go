package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/mobingilabs/pullr/pkg/comm"
	"github.com/mobingilabs/pullr/pkg/comm/rabbitmq"
	"github.com/mobingilabs/pullr/pkg/domain"
)

func main() {
	mqueue, err := rabbitmq.Dial("amqp://localhost:5672")
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	defer mqueue.Close()

	listener, err := mqueue.Listen(domain.BuildQueue)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	defer listener.Close()

	listenerCtx, cancelListen := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	log.Print("Waiting for jobs...")
	running := true

	go func() {
		<-c
		running = false
		cancelListen()
	}()

	for running {
		job, err := listener.Get(listenerCtx)
		if err != nil {
			if err != context.Canceled {
				log.Printf("[ERROR] Failed to get job from listener: %s", err)
			}
			continue
		}

		handleJob(job)
	}
}

func handleJob(job comm.Job) {
	body := job.Body()
	log.Printf("[INFO] Job body: %s", string(body))

	decoder := json.NewDecoder(bytes.NewReader(body))
	var buildJob domain.BuildImageJob
	if err := decoder.Decode(&buildJob); err != nil && err != io.EOF {
		log.Printf("[ERROR] Failed to parse job: %s", err)
		return
	}

	log.Printf("[INFO] Got new job: %+v", buildJob)
	//if err := job.Finish(); err != nil {
	//	log.Printf("[ERROR] Failed to mark job as finished: %s", err)
	//}
}
