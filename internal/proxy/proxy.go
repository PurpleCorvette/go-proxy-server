// Package proxy provides a gRPC proxy server for forwarding requests to backend services.
package proxy

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-proxy-server/internal/balancer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// ProxyServer defines an interface for forwarding requests to backend services.
type ProxyServer interface {
	StreamDirector(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error)
}

// gRPCProxyServer implements ProxyServer interface for gRPC.
type gRPCProxyServer struct {
	balancer balancer.Balancer
	logger   *logrus.Logger
}

// NewProxyServer initializes a new gRPCProxyServer with the provided balancer and logger.
func NewProxyServer(balancer balancer.Balancer, logger *logrus.Logger) ProxyServer {
	return &gRPCProxyServer{
		balancer: balancer,
		logger:   logger,
	}
}

// StreamDirector is a gRPC stream director for forwarding requests to backend services.
func (p *gRPCProxyServer) StreamDirector(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
	p.logger.WithFields(logrus.Fields{
		"method": fullMethodName,
	}).Info("Request received")

	conn := p.balancer.GetConn()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	} else {
		md = md.Copy()
	}
	if p, ok := peer.FromContext(ctx); ok {
		md.Set("x-forwarded-for", p.Addr.String())
	}
	newCtx := metadata.NewOutgoingContext(ctx, md)
	return newCtx, conn, nil
}
