package main

import (
	"sync"
)

const (
	restPort = 8080
	grpcPort = 50051
)

func main() {
	var wg sync.WaitGroup
	var tfGenCmd = []string{"./tfgen", "--host", "127.0.0.1:50051",
		"--out", config.outdir, "--rando"}

	// grpc server - pubsub connection with running python scripts
	cmds := NewCronMap()
	server := NewManager(grpcPort, cmds)

	// start running tensorflow script generator
	cmds.AddCmd("", func() []string {
		return tfGenCmd
	}, config.freq)

	// http server - client-side controller
	ctrl := NewController(restPort)
	ctrl.mgr = server
	ctrl.cmds = cmds

	wg.Add(2)
	go server.Start(&wg)
	go ctrl.Start(&wg)
	go cmds.Run()
	// runCmd(tfGenCmd[:len(tfGenCmd)-1])
	wg.Wait()
}
