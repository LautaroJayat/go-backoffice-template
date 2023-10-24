package propagation

import (
	"log"
	"testing"

	"github.com/lautarojayat/backoffice/products"
	users "github.com/lautarojayat/backoffice/users"
)

func testForProducts(t *testing.T, p *Publisher) {
	f, err := NewPublisherFunction(p, "", products.ProductOp{})
	if f != nil {
		t.Error("pubFun should be nil when channel is empty")
	}
	if err == nil {
		t.Error("err should not be nil when channel is empty")
	}
	f, err = NewPublisherFunction(p, "chan", products.ProductOp{})
	if f == nil {
		t.Error("pubFun should not be nil when channel is not empty and ProductOp can be converted to json")
	}
	if err != nil {
		t.Errorf("error should be nil when channel is not empty and ProductOp can be converted to json. got=%s", err)
	}
}

func testForCustomers(t *testing.T, p *Publisher) {
	f, err := NewPublisherFunction(p, "", users.UsersOp{})

	if f != nil {
		t.Error("pubFun should be nil when channel is empty")
	}
	if err == nil {
		t.Error("err should not be nil when channel is empty")
	}

	f, err = NewPublisherFunction(p, "chan", users.UsersOp{})
	if f == nil {
		t.Error("pubFun should not be nil when channel is not empty and UsersOp can be converted to json")
	}
	if err != nil {
		t.Errorf("error should  be nil when channel is not empty and UsersOp can be converted to json. got=%s", err)
	}
}

func TestNewPublisherFunction(t *testing.T) {
	p := &Publisher{l: log.Default()}
	testForCustomers(t, p)
	testForProducts(t, p)
}
