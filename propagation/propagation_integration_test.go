package propagation

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/lautarojayat/backoffice/config"
	"github.com/lautarojayat/backoffice/products"
	"github.com/redis/go-redis/v9"
)

func TestPropagationCapabilities(t *testing.T) {
	cfg := config.Redis{
		Addr:       "localhost:6379",
		PubTimeOut: 2,
	}

	subscriber := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		DB:   0})
	defer subscriber.Close()
	l := log.Default()

	publisher := NewPublisher(cfg, l)
	defer publisher.Stop()

	pubChannel := "products"
	pubFun, err := NewPublisherFunction[products.ProductOp](publisher, pubChannel, products.ProductOp{})

	if err != nil {
		t.Errorf("expected nil error, got=%q", err)
	}

	if pubFun == nil {
		t.Error("expected non nil pubFun")
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.PubTimeOut)*time.Second)
	defer cancel()

	expected := products.ProductOp{
		Op: "Create",
		Payload: products.Product{
			Name:  "p1",
			Price: 1,
		}}

	ctx2, cancel2 := context.WithTimeout(ctx, time.Duration(cfg.PubTimeOut)*time.Second)
	defer cancel2()
	sub := subscriber.Subscribe(ctx2, pubChannel)

	pubFun(ctx, expected)

	ctx3, cancel3 := context.WithTimeout(ctx, time.Duration(cfg.PubTimeOut)*time.Second)
	defer cancel3()

	msg, err := sub.ReceiveMessage(ctx3)
	if err != nil {
		t.Errorf("redis client should be able to receive message while testing. error=%q", err)
	}

	received := &products.ProductOp{}
	json.Unmarshal([]byte(msg.Payload), received)

	if received.Op != expected.Op {
		t.Errorf("received op must be %q, instead got %q", expected.Op, received.Op)
	}
	if received.Payload.Name != expected.Payload.Name {
		t.Errorf("received name must be %q, instead got %q", expected.Payload.Name, received.Payload.Name)
	}
	if received.Payload.Name != expected.Payload.Name {
		t.Errorf("received name must be %q, instead got %q", expected.Payload.Name, received.Payload.Name)
	}
	if received.Payload.Price != expected.Payload.Price {
		t.Errorf("received price must be %q, instead got %q", expected.Payload.Name, received.Payload.Name)
	}
	if received.Payload.CreatedAt.String() != expected.Payload.CreatedAt.String() {
		t.Errorf("received createdAt must be %q, instead got %q", expected.Payload.CreatedAt.String(), received.Payload.CreatedAt.String())
	}

}
