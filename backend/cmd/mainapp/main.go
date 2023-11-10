package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pazbear-backend/cmd/mainapp/config"
	"pazbear-backend/cmd/mainapp/controller"
	"pazbear-backend/cmd/mainapp/server"
	"syscall"
	"time"
)

func main() {

	log.Println("starting server...")
	cnf, err := config.AppConfig()
	if err != nil {
		log.Fatalf("Configuration is invalid!: %s", err.Error())
	}

	c, err := controller.NewController(cnf)
	if err != nil {
		panic(err)
	}

	r := c.NewRouter()
	svr := server.NewServer(r, fmt.Sprintf(":%d", cnf.Port))
	go func() {
		if err := svr.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error while running server: %s\n", err.Error())
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
	// stop server
	if err := svr.Stop(ctx); err != nil {
		log.Fatalf("failed to stop server: %v", err)
	}
}
