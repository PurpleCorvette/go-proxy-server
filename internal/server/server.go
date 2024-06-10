// Package server provides functionality to start the gRPC proxy server.
package server

import (
	"net"

	"github.com/mwitkow/grpc-proxy/proxy"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go-proxy-server/internal/auth"
	"go-proxy-server/internal/balancer"
	"go-proxy-server/internal/cache"
	"go-proxy-server/internal/config"
	"go-proxy-server/internal/metrics"
	internalProxy "go-proxy-server/internal/proxy"
)

// StartGRPCServer starts the gRPC proxy server with the provided configuration and logger.
func StartGRPCServer(cfg *config.Config, log *logrus.Logger) {
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create the load balancer
	balance, err := balancer.NewRoundRobinBalancer(cfg.Servers, log)
	if err != nil {
		log.Fatalf("Failed to create balancer: %v", err)
	}

	// Create the cache
	cacheStorage := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB, log)

	// Create the proxy server
	proxyServer := internalProxy.NewProxyServer(balance, log)

	// Set up the gRPC server with middleware
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			auth.NewJWTAuthenticator().AuthUnaryInterceptor(log),
			cacheStorage.CacheUnaryInterceptor(),
			metrics.NewPrometheusMetrics().MetricsUnaryInterceptor(log),
		),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(proxyServer.StreamDirector)),
	)

	reflection.Register(s)

	log.Infof("Starting gRPC proxy server on port: %s", cfg.Port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
