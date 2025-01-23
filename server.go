package main

import (
	"XcxcPan/Kafka"
	"XcxcPan/StorageGroup"
	"XcxcPan/etcd"
	"XcxcPan/router"
	"fmt"
	"log"
	"time"
)

func main() {
	r := router.Router()
	//初始哈希环，添加两个服务节点

	go func() {
		Kafka.StartConsumer_Del()
	}()
	go func() {
		StorageGroup.Server = StorageGroup.NewStorageServer("")
		cli, err := etcd.ClientInit()
		if err != nil {
			log.Fatalf("etcd client init failed: %v", err)
		}
		storagePeers, err := etcd.GetAllPeers(cli, "xcstorage")
		if err != nil {
			log.Fatalf("get storage nodes from etcd failed: %v", err)
		}
		fmt.Println("get storage nodes from etcd :", storagePeers)
		StorageGroup.Server.SetPeers(storagePeers...)

		go func() {
			time.Sleep(5 * time.Second)
			for {
				currentPeers, _ := etcd.GetAllPeers(cli, "xcstorage")
				closedPeers := GetClosedPeers(storagePeers, currentPeers)
				newPeers := GetNewPeers(storagePeers, currentPeers)
				if len(closedPeers) != 0 {
					StorageGroup.Server.DelPeers(closedPeers...)
					log.Println("Storage nodes have been closed :", closedPeers)
					log.Println("Current nodes :", currentPeers)
				}
				if len(newPeers) != 0 {
					StorageGroup.Server.AddPeers(newPeers...)
					log.Println("New storage nodes have been added :", newPeers)
					log.Println("Current nodes :", currentPeers)
				}
				storagePeers = currentPeers
				time.Sleep(5 * time.Second)

			}
		}()
	}()
	r.Run(":7090")

}
func GetClosedPeers(peers, newpeers []string) []string {
	var closedPeers []string
	for _, peer := range peers {
		if !IsContain(newpeers, peer) {
			closedPeers = append(closedPeers, peer)
		}
	}
	return closedPeers
}
func GetNewPeers(peers, newpeers []string) []string {
	var newPeers []string
	for _, peer := range newpeers {
		if !IsContain(peers, peer) {
			newPeers = append(newPeers, peer)
		}
	}
	return newPeers
}

func IsContain(peers []string, peer string) bool {
	for _, p := range peers {
		if p == peer {
			return true
		}
	}
	return false
}
