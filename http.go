package main

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/raft"
)

type HttpServer struct {
	raft raft.Raft
}

func (h *HttpServer) Start(httpAddr string) {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(httpAddr, nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
