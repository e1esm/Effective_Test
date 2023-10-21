package app

import (
	"context"
	"github.com/e1esm/Effective_Test/internal/server"
	"log"
	"os/signal"
	"syscall"
)

func Run() {
	srv := server.NewHttpServer()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()
	<-ctx.Done()
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("error occurred while shutting down application: %v", err.Error())
	}
}
