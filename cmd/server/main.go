package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpApi "github.com/lautarojayat/backoffice/api/http"
	"github.com/lautarojayat/backoffice/config"
	"github.com/lautarojayat/backoffice/logger"
	database "github.com/lautarojayat/backoffice/persistence/db"
	"github.com/lautarojayat/backoffice/products"
	"github.com/lautarojayat/backoffice/propagation"
	users "github.com/lautarojayat/backoffice/users"

	"github.com/lautarojayat/backoffice/persistence/files"
	"github.com/lautarojayat/backoffice/server"
)

func listenAndServe(s *http.Server, notifyEnd chan struct{}) {
	err := s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Printf("error while serving http. error=%q", err)
		notifyEnd <- struct{}{}
	}

}

func main() {
	l := logger.New()
	cfgFile, err := files.OpenFile("config/default.yaml")
	l.Println("about to read configs")

	if err != nil {
		l.Fatalf("fatal: couldn't read cfg file. error=%q", err)
	}

	l.Println("about to generate configs object")
	cfg, err := config.FromYAML(cfgFile)

	l.SetPrefix(cfg.AppName)

	if err != nil {
		l.Fatalf("fatal: couldn't create config object. error=%q", err)
	}

	l.Println("about to connect with DB")

	db, err := database.NewConnection(cfg.DB)
	if err != nil {
		l.Fatalf("fatal: couldn't connect to db. error=%q", err)
	}

	err = database.RunMigrations(l, db)
	if err != nil {
		l.Fatalf("fatal: couldn't perform all migrations. error=%q", err)
	}

	pub := propagation.NewPublisher(cfg.Propagation.Redis, l)

	usersPubFun, err := propagation.NewPublisherFunction(
		pub,
		cfg.Propagation.Channels.Customers,
		users.UsersOp{},
	)
	if err != nil {
		l.Fatal(err)
	}
	usersRepo := users.NewRepo(db, usersPubFun)

	productPubFun, err := propagation.NewPublisherFunction(
		pub,
		cfg.Propagation.Channels.Products,
		products.ProductOp{},
	)
	if err != nil {
		l.Fatal(err)
	}
	productsRepo := products.NewRepo(db, productPubFun)

	l.Println("about to generate http endpoints")
	mux := httpApi.MakeHTTPEndpoints(l, usersRepo, productsRepo)

	l.Println("about to listen and serve")
	s := server.NewServer(cfg.HTTP, mux)
	notifyEnd := make(chan struct{})
	go listenAndServe(s, notifyEnd)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigs:
		l.Printf("received signal %q", sig)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15)*time.Second)
		defer cancel()
		err = s.Shutdown(ctx)
		if err != nil {
			l.Fatalf("fatal: error while shutting down. error=%q", err)
		}
		pub.Stop()

	case <-notifyEnd:
		pub.Stop()
	}

	l.Println("Terminating process")
}
