package main

import (
	"cluster-balance/comm"
	"cluster-balance/master"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("Master START !!!")
	config := &comm.Config{
		ClusterID:  "sunwei_test",
		EetcdHosts: []string{"http://127.0.0.1:2379"},
	}

	master, err := master.NewMaster(config, "test_Master_1")
	if err != nil {
		log.Error("master.NewMaster, err:", master)
		return
	}

	go master.Run()

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

	master.Stop()

	log.Info("Node_Master Stop")
}
