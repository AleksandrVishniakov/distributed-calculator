package userscache

import (
	"container/list"
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/domain/models"
)

type UserProvider interface {
	User(
		ctx context.Context,
		login string,
	) (user *models.User, err error)
}

type cacheItem struct {
	user     *models.User
	issuedAt time.Time
}

type Cache struct {
	log *slog.Logger

	userProvider UserProvider

	mu          *sync.RWMutex
	listLock    *list.List
	storageLock map[string]*list.Element

	maxSize int
	ttl     time.Duration
}

func New(
	log *slog.Logger,
	userProvider UserProvider,
	maxSize int,
	ttl time.Duration,
) *Cache {
	return &Cache{
		log:          log,
		userProvider: userProvider,
		mu:           &sync.RWMutex{},
		storageLock:  make(map[string]*list.Element),
		maxSize:      maxSize,
		ttl:          ttl,
		listLock:     list.New(),
	}
}

func (c *Cache) User(ctx context.Context, login string) (user *models.User, err error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	const src = "Cache.User"

	log := c.log.With(
		slog.String("src", src),
		slog.String("login", login),
	)

	log.Debug("trying to get user...")

	el, exists := c.element(login)

	var item = &cacheItem{issuedAt: time.Now()}

	if exists {
		var ok bool
		item, ok = el.Value.(*cacheItem)
		if !ok {
			return nil, errors.New("list.Element type conversion failed")
		}
	}

	var isItemExpired = time.Now().After(item.issuedAt.Add(c.ttl))
	if !exists || isItemExpired {
		if !exists {
			log.Debug("cache miss")
		} else if isItemExpired {
			log.Debug("item expired")
			c.invalidate(login)
		}

		user, err := c.userProvider.User(ctx, login)
		if err != nil {
			return nil, err
		}

		err = c.Set(ctx, user)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	c.markAsUsed(el)

	return item.user, nil
}

func (c *Cache) Set(ctx context.Context, user *models.User) (err error) {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	const src = "Cache.Set"

	log := c.log.With(
		slog.String("src", src),
		slog.String("login", user.Login),
	)

	log.Debug("caching user...")

	item := &cacheItem{user: user, issuedAt: time.Now()}

	if c.length() >= c.maxSize {
		lastEl := c.leastCommonlyUsed()

		item, ok := lastEl.Value.(*cacheItem)
		if !ok {
			return errors.New("list.Element type conversion failed")
		}

		log.Debug("storage is overflowed. invalidate " + item.user.Login)

		c.invalidate(item.user.Login)
	}

	c.save(user.Login, item)

	log.Debug("user cached")

	return nil
}

func (c *Cache) save(key string, item *cacheItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el := c.listLock.PushFront(item)
	c.storageLock[key] = el
}

func (c *Cache) element(login string) (*list.Element, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	el, ok := c.storageLock[login]
	return el, ok
}

func (c *Cache) leastCommonlyUsed() *list.Element {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.listLock.Back()
}

func (c *Cache) invalidate(login string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.storageLock[login]; ok {
		c.listLock.Remove(el)
		delete(c.storageLock, login)
	}
}

func (c *Cache) length() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.storageLock)
}

func (c *Cache) markAsUsed(el *list.Element) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.listLock.MoveToFront(el)
}
