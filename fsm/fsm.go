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
	DataBase database
}

func NewFsm() *Fsm {
	fsm := &Fsm{
		DataBase: NewDatabase(),
	}
	return fsm
}

func (f *Fsm) Apply(l *raft.Log) interface{} {
	fmt.Println("apply data:", string(l.Data))
	data := strings.Split(string(l.Data), ",")
	op := data[0]
	if op == "set" {
		key := data[1]
		value := data[2]
		f.DataBase.Set(key, value)
	}

	return nil
}

func (f *Fsm) Snapshot() (raft.FSMSnapshot, error) {
	return &f.DataBase, nil
}

func (f *Fsm) Restore(io.ReadCloser) error {
	return nil
}

type database struct {
	Data map[string]string
	mu   sync.Mutex
}

func NewDatabase() database {
	return database{
		Data: make(map[string]string),
	}
}

func (d *database) Get(key string) string {
	d.mu.Lock()
	value := d.Data[key]
	d.mu.Unlock()
	return value
}

func (d *database) Set(key, value string) {
	d.mu.Lock()
	d.Data[key] = value
	d.mu.Unlock()
}

func (d *database) Persist(sink raft.SnapshotSink) error {
	d.mu.Lock()
	data, err := json.Marshal(d.Data)
	d.mu.Unlock()
	if err != nil {
		return err
	}
	sink.Write(data)
	sink.Close()
	return nil
}

func (d *database) Release() {}
