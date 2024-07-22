package app

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"manga/internal/http"
	"manga/internal/message_broker/broker"
	"manga/internal/store/postgres"

	"os"
	"os/signal"
	"syscall"
)

func run() {
	ctx, cancel := context.WithCancel(context.Background())
	go CatchTermination(cancel)

	dbURL := "postgres://postgres:postgres@localhost:5432/postgres"
	store := postgres.NewDB()
	if err := store.Connect(dbURL); err != nil {
		panic(err)
	}
	defer store.Close()

	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}

	brokers := []string{"localhost:"}
	broker := broker.NewBroker(brokers,cache,"Name")
	if err := broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer broker.Close()

	srv := http.NewServer(
		ctx,
		http.WithAddress(":8082"),
		http.WithStore(store),
		http.WithCache(cache),
		http.WithBroker(broker),
	)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGraceFulTarmination()
}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Print("[WARN] caught termination signal")
	cancel()
}
