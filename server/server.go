package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	gpb "github.com/mingkaic/go_tenncor/graphmgr"
	spb "github.com/mingkaic/go_tenncor/serial"
	log "github.com/sirupsen/logrus"
)

type GraphMgr struct {
	cmds CronMap
}

func (s *GraphMgr) PostGraph(ctx context.Context, in *gpb.GraphCreated) (*gpb.Nothing, error) {
	log.Infof("generated graph id: %s", in.Gid)
	s.cmds.AddCmd(in.Gid, func() []string {
		dataID := uuid.New().String()
		out := []string{"python", config.outdir + "/" + in.Gid + ".py", dataID}
		// store dataID
		log.Infof("store data id %s", dataID)
		return out
	}, 15*time.Second)
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
	manager  *GraphMgr
	listener net.Listener
}

func NewManager(port int, cmds CronMap) *GRPCServer {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Graphmgr failed to listen: %v", err)
	}

	server := grpc.NewServer()
	manager := &GraphMgr{cmds: cmds}
	gpb.RegisterGraphmgrServer(server, manager)
	return &GRPCServer{
		server,
		manager,
		listener,
	}
}

// Start ... runs grpc server at port specified
func (server *GRPCServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info("Graphmgr server running")
	err := server.Serve(server.listener)
	if err != nil && err != grpc.ErrServerStopped {
		log.Fatalf("Graphmgr server error: %v", err)
	}
	log.Info("Graphmgr server stopped")
}

func (server *GRPCServer) Stop() {
	server.Server.Stop()
}
