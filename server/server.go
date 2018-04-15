package main

import (
	"sync"
	"fmt"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	gpb "github.com/mingkaic/go_tenncor/graphmgr"
	spb "github.com/mingkaic/go_tenncor/serial"
	log "github.com/sirupsen/logrus"
)

type GraphMgr struct{}

func (s *GraphMgr) PostGraph(ctx context.Context, in *gpb.GraphCreated) (*gpb.Nothing, error) {
	log.Println("generated graph id:", in.Gid)
	return &gpb.Nothing{}, nil
}

func (s *GraphMgr) GetGraphList(ctx context.Context, in *gpb.GraphListReq) (*gpb.GraphList, error) {
	return &gpb.GraphList{
		Gids: []string{"sample"},
	}, nil
}

func (s *GraphMgr) GetGraphPb(ctx context.Context, in *gpb.GraphReq) (*spb.GraphPb, error) {
	createOrder := []string{"A", "B", "C"}
	nodeMap := make(map[string]*spb.NodePb)
	nodeMap["A"] = &spb.NodePb{}
	nodeMap["B"] = &spb.NodePb{}
	nodeMap["C"] = &spb.NodePb{}
	return &spb.GraphPb{
		Gid:         "sample",
		CreateOrder: createOrder,
		NodeMap:     nodeMap,
	}, nil
}

func (s *GraphMgr) GetData(in *gpb.DataReq, stream gpb.Graphmgr_GetDataServer) error {
	return nil
}

func (s *GraphMgr) GetTestData(in *gpb.DataReq, stream gpb.Graphmgr_GetTestDataServer) error {
	return nil
}

type GRPCServer struct {
	*grpc.Server
	listener net.Listener
}

func NewManager(port int) *GRPCServer {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Graphmgr failed to listen: %v", err)
	}

	server := grpc.NewServer()
	gpb.RegisterGraphmgrServer(server, &GraphMgr{})
	return &GRPCServer{
		server,
		listener,
	}
}

// RunGRPCServer runs grpc server at port specified
func (server *GRPCServer) Start(wg sync.WaitGroup) {
	defer wg.Done()
	log.Info("Graphmgr server running")
	err := server.Serve(server.listener)
	if err != nil && err != grpc.ErrServerStopped {
		log.Fatalf("Graphmgr server error: %v", err)
	}
	log.Info("Graphmgr server stopped")
}
