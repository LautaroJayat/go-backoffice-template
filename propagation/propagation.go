package propagation

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lautarojayat/e_shop/config"
	"github.com/lautarojayat/e_shop/products"
	users "github.com/lautarojayat/e_shop/users"
	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	available  bool
	lock       *sync.RWMutex
	l          *log.Logger
	rdb        *redis.Client
	pubTimeOut time.Duration
}

type Propagable interface {
	users.UsersOp | products.ProductOp
}

func NewPublisher(cfg config.Redis, l *log.Logger) *Publisher {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		Username: cfg.User,
		DB:       0, // use default DB
	})

	return &Publisher{
		available:  true,
		lock:       &sync.RWMutex{},
		l:          l,
		rdb:        rdb,
		pubTimeOut: time.Millisecond * time.Duration(cfg.PubTimeOut),
	}
}

func (p *Publisher) Stop() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.available = false
	p.rdb.Close()
}

func NewPublisherFunction[T Propagable](p *Publisher, channel string, entity T) (func(context.Context, T), error) {
	if channel == "" {
		return nil, fmt.Errorf("could not create publisher function due to empty string used as channel name")
	}

	_, err := json.Marshal(&entity)

	if err != nil {

		return nil, fmt.Errorf("could not create publisher function for %T. error=%q", entity, err)
	}

	return func(ctx context.Context, payload T) {
		p.lock.RLock()
		defer p.lock.RUnlock()

		if !p.available {
			p.l.Printf("attempted to publish in channel %q while Publisher was unavailable", channel)
			return
		}

		s, err := json.Marshal(payload)

		if err != nil {
			p.l.Printf("could not publish to channel %q due to bad payload interface. error=%q", channel, err)
		}

		err = p.rdb.Publish(ctx, channel, s).Err()

		if err != nil {
			p.l.Printf("could not publish to channel %q due to client/server error. error=%q", channel, err)
		}
	}, nil
}
