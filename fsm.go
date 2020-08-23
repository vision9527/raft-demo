package main

import (
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/raft"
)

type Fsm struct {
	mu   sync.Mutex
	data map[string]string
}

func (f *Fsm) Apply(l *raft.Log) interface{} {
	fmt.Println("apply data:", string(l.Data))
	return nil
}

func (f *Fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return nil, nil
}

func (f *Fsm) Restore(io.ReadCloser) error {
	return nil
}
