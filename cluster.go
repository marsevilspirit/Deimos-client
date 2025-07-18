package deimosclient

import (
	"log/slog"
	"math/rand"
	"sync"
)

type Cluster struct {
	Leader    string   `json:"leader"`
	Endpoints []string `json:"endpoints"`
	picked    int
	mu        sync.RWMutex
}

func NewCluster(endpoints []string) *Cluster {
	// if an empty slice was sent in then just assume HTTP 4001 on localhost
	if len(endpoints) == 0 {
		endpoints = []string{"http://127.0.0.1:4001"}
	}

	endpoints = shuffleStringSlice(endpoints)
	slog.Debug("Shuffle cluster", "machines", endpoints)
	// default leader and machines
	return &Cluster{
		Leader:    "",
		Endpoints: endpoints,
		picked:    rand.Intn(len(endpoints)),
	}
}

func (cl *Cluster) pick() string {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.Endpoints[cl.picked]
}
