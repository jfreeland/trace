package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jfreeland/trace/api"
	"github.com/jfreeland/trace/storage"
	"github.com/jfreeland/trace/tracer"
)

func main() {
	// do something with a config
	port := 50000

	router := api.New()
	db := storage.NewInMemory()
	api.AddRoutes(router, db)

	tracer := tracer.NewAeden(db)
	go tracer.Run("google.com")

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           router,
		IdleTimeout:       time.Second * 60,
		ReadTimeout:       time.Second * 10,
		WriteTimeout:      time.Second * 10,
		ReadHeaderTimeout: 1 * time.Second,
		TLSNextProto:      map[string]func(*http.Server, *tls.Conn, http.Handler){}, // disable HTTP/2
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Errorf("listen: %v", err)
		}
	}()

	go func() {
		time.Sleep(15 * time.Second)
		tracer.Stop("google.com")
	}()

	log.Infof("serving on: %d", port)

	// Wait for interrupt signal to gracefully shutdown the server with a
	// timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	log.Infof("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Debugf("server shutdown failed: %v", err)
	}
	log.Infof("shutdown complete")
}
