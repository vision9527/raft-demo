package fsm

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/hashicorp/raft"
)

type Fsm struct {
	mu   sync.Mutex
	Data database
}

func (f *Fsm) Apply(l *raft.Log) interface{} {
	fmt.Println("apply data:", string(l.Data))
	data := strings.Split(string(l.Data), ",")
	op := data[0]
	f.mu.Lock()
	if op == "set" {
		key := data[1]
		value := data[2]
		f.Data[key] = value
	}
	f.mu.Unlock()

	return nil
}

func (f *Fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.Data, nil
}

func (f *Fsm) Restore(io.ReadCloser) error {
	return nil
}

type database map[string]string

func (d database) Persist(sink raft.SnapshotSink) error {
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	sink.Write(data)
	sink.Close()
	return nil
}

func (f database) Release() {}
