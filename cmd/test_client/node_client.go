package main

import (
	"cluster-balance/api"
	"cluster-balance/client"
	"cluster-balance/comm"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

type Worker struct {
}

func (w *Worker) HanderDoPut(key, value []byte) error {
	log.Info("HanderDoPut, key:", string(key), " value:", string(value))
	return nil
}

func (w *Worker) HanderDoDelete(key, value []byte) error {
	log.Info("HanderDoDelete, key:", string(key), " value:", string(value))
	return nil
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("Client START")
	config := &comm.Config{
		ClusterID:  "sunwei_test",
		EetcdHosts: []string{"http://127.0.0.1:2379"},
	}

	woker := &Worker{}
	node := &api.NodeSpec{
		ClusterId: config.ClusterID,
		Id:        "test_client1",
		Quotas:    map[string]uint64{"channel": 10},
		Labels:    map[string]string{"video": "face"},
	}

	c, err := client.NewClient(config, node, woker)
	if err != nil {
		log.Error("client.NewClient failed, err:", err)
		return
	}

	go func() {
		if err := c.Run(); err != nil {
			log.Fatal("client exited: ", err)
		}
	}()

	log.Info("Node_Client Start")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

MainLoop:
	for {
		select {
		case <-sc:
			break MainLoop
		}
	}

	c.Stop()

	log.Info("Node_Client Stop")
}
