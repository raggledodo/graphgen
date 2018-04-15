package main

import (
	"fmt"
	"sync"
)

const (
	restPort = 8080
	grpcPort = 50051
)

func main() {
	var wg sync.WaitGroup

	// grpc server - pubsub connection with running python scripts
	server := NewManager(grpcPort)

	// start running tensorflow script generator
	job := NewCron(config.schedule)

	// http server - client-side controller
	ctrl := NewController(restPort)
	ctrl.mgr = server
	ctrl.job = job

	wg.Add(2)
	go server.Start(wg)
	go job.Run()
	go ctrl.Start(wg)
	wg.Wait()
}
