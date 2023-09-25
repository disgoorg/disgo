package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

var (
	token = os.Getenv("disgo_token")
)

func main() {
	slog.Info("starting example...")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds|gateway.IntentGuildMessages|gateway.IntentDirectMessages)),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds|cache.FlagMessages|cache.FlagMembers),
			cache.WithMemberCache(cache.NewMemberCache(newGroupedCache[discord.Member]())),
		),
	)
	if err != nil {
		slog.Error("error while building bot", slog.Any("err", err))
		return
	}

	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		client.Close(closeCtx)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = client.OpenGateway(ctx); err != nil {
		slog.Error("error while connecting to gateway", slog.Any("err", err))
	}

	slog.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

var _ cache.GroupedCache[discord.Member] = (*groupedCache[discord.Member])(nil)

func newGroupedCache[T any]() cache.GroupedCache[T] {
	return &groupedCache[T]{
		cache: make(map[snowflake.ID]map[snowflake.ID]T),
	}
}

type groupedCache[T any] struct {
	cache map[snowflake.ID]map[snowflake.ID]T
	mu    sync.Mutex
}

func (g *groupedCache[T]) Get(groupID snowflake.ID, id snowflake.ID) (T, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	groupEntities, ok := g.cache[groupID]
	if !ok {
		var entity T
		return entity, false
	}

	entity, ok := groupEntities[id]
	return entity, ok
}

func (g *groupedCache[T]) Put(groupID snowflake.ID, id snowflake.ID, entity T) {
	g.mu.Lock()
	defer g.mu.Unlock()

	groupEntities, ok := g.cache[groupID]
	if !ok {
		groupEntities = make(map[snowflake.ID]T)
		g.cache[groupID] = groupEntities
	}

	groupEntities[id] = entity
}

func (g *groupedCache[T]) Remove(groupID snowflake.ID, id snowflake.ID) (T, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	groupEntities, ok := g.cache[groupID]
	if !ok {
		var entity T
		return entity, false
	}

	entity, ok := groupEntities[id]
	if !ok {
		return entity, false
	}

	delete(groupEntities, id)
	return entity, true
}

func (g *groupedCache[T]) GroupRemove(groupID snowflake.ID) {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.cache, groupID)
}

func (g *groupedCache[T]) RemoveIf(filterFunc cache.GroupedFilterFunc[T]) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for groupID, groupEntities := range g.cache {
		for id, entity := range groupEntities {
			if filterFunc(groupID, entity) {
				delete(groupEntities, id)
			}
		}
	}
}

func (g *groupedCache[T]) GroupRemoveIf(groupID snowflake.ID, filterFunc cache.GroupedFilterFunc[T]) {
	g.mu.Lock()
	defer g.mu.Unlock()

	groupEntities, ok := g.cache[groupID]
	if !ok {
		return
	}

	for id, entity := range groupEntities {
		if filterFunc(groupID, entity) {
			delete(groupEntities, id)
		}
	}
}

func (g *groupedCache[T]) Len() int {
	g.mu.Lock()
	defer g.mu.Unlock()

	var length int
	for _, groupEntities := range g.cache {
		length += len(groupEntities)
	}
	return length
}

func (g *groupedCache[T]) GroupLen(groupID snowflake.ID) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	groupEntities, ok := g.cache[groupID]
	if !ok {
		return 0
	}

	return len(groupEntities)
}

func (g *groupedCache[T]) ForEach(f func(groupID snowflake.ID, entity T)) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for groupID, groupEntities := range g.cache {
		for _, entity := range groupEntities {
			f(groupID, entity)
		}
	}
}

func (g *groupedCache[T]) GroupForEach(groupID snowflake.ID, forEachFunc func(entity T)) {
	g.mu.Lock()
	defer g.mu.Unlock()

	groupEntities, ok := g.cache[groupID]
	if !ok {
		return
	}

	for _, entity := range groupEntities {
		forEachFunc(entity)
	}
}
