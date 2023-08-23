package main

import (
	"context"
	"sort"
	"strings"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"
	"golang.org/x/exp/slices"
)

var ( 
	cacheHits = metrics.NewCounter( 
		"search_cache_hits", 
		"Number of Search cache hits", 
	) 
	cacheMisses = metrics.NewCounter( 
		"search_cache_misses", 
		"Number of Search cache misses", 
	) 
) 

type Searcher interface {
	Search(context.Context, string) ([]string, error)
}

type searcher struct {
	weaver.Implements[Searcher]
	cache weaver.Ref[Cache] 
}

func (s *searcher) Search(ctx context.Context, query string) ([]string, error) {
	s.Logger(ctx).Debug("Search", "query", query)

	if emojis, err := s.cache.Get().Get(ctx, query); err != nil {
		s.Logger(ctx).Error("cache.Get", "query", query, "err", err)
	} else if len(emojis) > 0 {
		cacheHits.Add(1) 
		return emojis, nil
	} else {
		cacheMisses.Add(1)
	}

	words := strings.Fields(strings.ToLower(query))
	var results []string
	for emoji, labels := range emojis {
		if matches(labels, words) {
			results = append(results, emoji)
		}
	}
	sort.Strings(results)

	if err := s.cache.Get().Put(ctx, query, results); err != nil { 
		s.Logger(ctx).Error("cache.Put", "query", query, "err", err) 
	} 

	return results, nil
}

func matches(labels, words []string) bool {
	for _, word := range words {
		if !slices.Contains(labels, word) {
			return false
		}
	}
	return true
}