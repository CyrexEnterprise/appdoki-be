package main

import (
	"appdoki-be/app"
	"appdoki-be/config"
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conf := config.NewConfig()

	db := prepareDatabase(&conf.Database)

	application := app.NewApplication(conf, db)

	srv := &http.Server{
		Addr:         conf.Server.Address,
		Handler:      application.Routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()
	log.Infof("Server started on %s", srv.Addr)

	<-done
	log.Info("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}
	log.Info("Server exited gracefully")
}

func prepareDatabase(conf *config.DatabaseConfig) *sqlx.DB {
	db, err := sqlx.Connect("postgres", conf.URI)
	if err != nil {
		log.Fatalln(err)
	}

	runMigrations(db.DB, conf)

	return db
}
