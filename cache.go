package main

import (
	"context"
	"sync"

	"github.com/ServiceWeaver/weaver"
)

// Cache caches emoji query results.
type Cache interface {
	// Get returns the cached emojis produced by the provided query. On cache
	// miss, Get returns nil, nil.
	Get(context.Context, string) ([]string, error)

	// Put stores a query and its corresponding emojis in the cache.
	Put(context.Context, string, []string) error
}

// cache implements the Cache component.
type cache struct {
	weaver.Implements[Cache]

	mu     sync.Mutex
	emojis map[string][]string
}

func (c *cache) Init(context.Context) error {
	c.emojis = map[string][]string{}
	return nil
}

func (c *cache) Get(ctx context.Context, query string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Logger(ctx).Debug("Get", "query", query)
	return c.emojis[query], nil
}

func (c *cache) Put(ctx context.Context, query string, emojis []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Logger(ctx).Debug("Put", "query", query)
	c.emojis[query] = emojis
	return nil
}