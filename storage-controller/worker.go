package main

import (
	"fmt"
	"sync"
)

// ReplicationData represent data to replicate between servers
type ReplicationData struct {
	Key string
	Val string
}

var numWorkers int = 2

func startReplicationWorkers(dataCh <-chan *ReplicationData, stopCh <-chan struct{}, wg *sync.WaitGroup) {
	for i := 0; i < numWorkers; i++ {
		go replicationWorker(dataCh, stopCh, wg)
	}
}

func replicationWorker(dataCh <-chan *ReplicationData, stopCh <-chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	fmt.Println("starting worker.")
	for {
		select {
		case data := <-dataCh:
			broadcast(data)
		case <-stopCh:
			fmt.Println("Stopping worker.")
			return
		}
	}
}

func broadcast(data *ReplicationData) {
	// TODO
	fmt.Printf("Replicating %v\n", data)
}
