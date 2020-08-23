package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	httpAddr string
	raftAddr string
	raftId   string
	raftDir  string
	joinAddr string
)

func init() {
	flag.StringVar(&httpAddr, "http_addr", "", "http listen addr")
	flag.StringVar(&raftAddr, "raft_addr", "", "raft listen addr")
	flag.StringVar(&raftId, "raft_id", "", "raft id")
	flag.StringVar(&joinAddr, "join_addr", "", "join addr")
}

func main() {
	fmt.Println("hello raft")
	flag.Parse()
	// 初始化配置
	if httpAddr == "" || raftAddr == "" || raftId == "" {
		fmt.Println("config error")
		os.Exit(1)
	}
	raftDir := "./node" + raftId
	os.MkdirAll(raftDir, 0700)

	// 初始化raft
	myRaft, err := NewMyRaft(raftAddr, raftId, raftDir)
	if err != nil {
		fmt.Println("NewMyRaft error")
		os.Exit(1)
	}

	// 启动raft/加入raft
	if joinAddr == "" {
		Bootstrap(myRaft)
	} else {
		// join
	}

	// 启动http server
	httpServer := new(HttpServer)
	httpServer.Start(httpAddr)

	// 退出http server
	fmt.Println("show kv http server")
}
