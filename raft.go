package main

import (
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

func NewMyRaft(raftAddr, raftId, raftDir string) (*raft.Raft, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(raftId)
	addr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(raftAddr, addr, 2, 5*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}
	snapshots, err := raft.NewFileSnapshotStore(raftDir, 2, os.Stderr)
	if err != nil {
		return nil, err
	}
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-log.db"))
	if err != nil {
		return nil, err
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-stable.db"))
	if err != nil {
		return nil, err
	}
	fsm := new(Fsm)

	rf, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return nil, err
	}

	return rf, nil
}

func Bootstrap(rf *raft.Raft) {
	config := rf.GetConfiguration().Configuration()
	rf.BootstrapCluster(config)
}
