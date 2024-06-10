// Package cache provides a gRPC unary interceptor for caching responses using Redis.
package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Cache defines an interface for caching responses.
type Cache interface {
	CacheUnaryInterceptor() grpc.UnaryServerInterceptor
}

// RedisCache implements Cache interface using Redis.
type RedisCache struct {
	client *redis.Client
	logger *logrus.Logger
}

// NewRedisCache initializes a new RedisCache with the provided Redis settings and logger.
func NewRedisCache(add, password string, db int, logger *logrus.Logger) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     add,
		Password: password,
		DB:       db,
	})
	return &RedisCache{client: rdb, logger: logger}
}

// CacheUnaryInterceptor returns a gRPC unary interceptor for caching responses.
func (c *RedisCache) CacheUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Try to get the cached response
		cachedResponse, err := c.client.Get(ctx, info.FullMethod).Result()
		if err == nil {
			c.logger.Infof("Cache hit for method: %s", info.FullMethod)
			return cachedResponse, nil
		}

		// Call the handler and cache the response
		resp, err := handler(ctx, req)
		if err == nil {
			c.client.Set(ctx, info.FullMethod, resp.(string), 5*time.Minute)
			c.logger.Infof("Cache miss for method: %s. Data cached", info.FullMethod)
		}
		return resp, err
	}
}
