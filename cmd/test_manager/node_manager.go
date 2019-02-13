package main

import (
	pb "cluster-balance/api"
	"cluster-balance/api/example"
	"cluster-balance/comm"
	"cluster-balance/master"
	"encoding/json"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)


type exampleServer struct {
	m *master.Master
}

type exampleConfig struct {
	etcd         comm.Config `json:"etcd"`
	GrpcEndpoint string      `json:"grpc_endpoint"`
	HttpEndpoint string      `json:"http_endpoint"`
	ManagerId    string      `json:"manager_id"`
}

func LoadConfig(fn string) (*exampleConfig, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	var c exampleConfig
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *exampleServer) AddResource(ctx context.Context, request *pb.ResourceSpec) (*example.ResourceResponse, error){
	if request.GetId() == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "empty resource id")
	}

	ns, err := s.m.AddResource(request)

	if err != nil {
		return nil, err
	}
	return &example.ResourceResponse{NodeId: ns.GetId()}, nil
}

func (s *exampleServer) RemoveResource(ctx context.Context, request *pb.ResourceSpec) (*example.ResourceResponse, error){
	if request.GetId() == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "empty resource id")
	}

	err := s.m.RemoveResource(request)
	if err != nil {
		return nil, err
	}
	return &example.ResourceResponse{}, nil
}

func (s *exampleServer) GetResource(ctx context.Context, request *pb.ResourceSpec) (*pb.ResourceSpec, error){
	return nil, nil
}

func (s *exampleServer) ListResources(ctx context.Context, request *example.ListResourcesRequest) (*example.ListResourcesResponse, error){
	return nil, nil
}

func (s *exampleServer) ListNodes(ctx context.Context, request *example.ListNodesRequest) (*example.ListNodesResponse, error){
	return nil, nil
}

func (s *exampleServer) GetNode(ctx context.Context, request *pb.NodeSpec) (*pb.NodeSpec, error){
	return nil, nil
}

var (
	configFile = flag.String("config", "./config.json", "cluster config file")
	verbose    = flag.Bool("verbose", false, "verbose output")
)

func main() {
	flag.Parse()
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Master START !!!")
	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatal("Load config file failed, err:", err)
	}

	master, err := master.NewMaster(&config.etcd, config.ManagerId)
	if err != nil {
		log.Error("master.NewMaster, err:", master)
		return
	}

	go master.Run()

	lis, err := net.Listen("tcp", config.GrpcEndpoint)
	if err != nil {
		log.Fatal("failed to listen, err:", err)
	}

	var serverOptions []grpc.ServerOption
	grpcServer := grpc.NewServer(serverOptions...)
	es := &exampleServer{
		m: master,
	}

	//start grpc server
	example.RegisterManagerServiceServer(grpcServer, es)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Warn("error: ", err)
		}
	}()

	//start http server
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	dialOptions := []grpc.DialOption{grpc.WithInsecure()}
	if err := example.RegisterManagerServiceHandlerFromEndpoint(ctx, mux, config.GrpcEndpoint, dialOptions); err != nil {
		log.Fatal("failed to register: ", err)
	}
	http.Handle("/", mux)
	httpServer := http.Server{Addr: config.HttpEndpoint}
	log.Info("gRPC http proxy listening at ", config.HttpEndpoint)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Info("http server closed: ", err)
		}
	}()

	log.Info("Node_Master Start")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

MainLoop:
	for {
		select {
		case <-sc:
			break MainLoop
		}
	}

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Warn(err)
	}
	grpcServer.GracefulStop()
	master.Stop()

	log.Info("Node_Master Stop")
}
