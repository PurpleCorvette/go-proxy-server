// Package balancer provides a load balancer for gRPC connections using a round-robin algorithm.
package balancer

import (
	"context"
	"golang.org/x/exp/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Balancer defines an interface for load balancing gRPC connections.
type Balancer interface {
	GetConn() *grpc.ClientConn
}

// RoundRobinBalancer implements Balancer interface using a round-robin algorithm.
type RoundRobinBalancer struct {
	servers []*grpc.ClientConn
	mu      sync.Mutex
	index   int
	logger  *logrus.Logger
}

// NewRoundRobinBalancer initializes a new RoundRobinBalancer with the provided server URLs and logger.
func NewRoundRobinBalancer(serverURLs []string, logger *logrus.Logger) (*RoundRobinBalancer, error) {
	servers := make([]*grpc.ClientConn, len(serverURLs))
	for i, serverURL := range serverURLs {
		conn, err := grpc.DialContext(context.Background(), serverURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		servers[i] = conn
	}
	return &RoundRobinBalancer{servers: servers, index: 0, logger: logger}, nil
}

// GetConn returns the next gRPC client connection using the round-robin algorithm.
func (b *RoundRobinBalancer) GetConn() *grpc.ClientConn {
	b.mu.Lock()
	defer b.mu.Unlock()
	server := b.servers[b.index]
	b.index = (b.index + 1) % len(b.servers)
	return server
}

// init initializes the seed for the random number generator.
func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}
