package main

import (
	"context"
	"dnd-backend-go/api/env"
	"dnd-backend-go/internal/client"
	"dnd-backend-go/internal/commands"
	"dnd-backend-go/internal/handlers"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func main() {
	var router *chi.Mux = chi.NewRouter()
	slog.SetLogLoggerLevel(slog.LevelDebug)

	var APP_PORT string
	APP_PORT = env.APP_PORT.GetValue()
	if APP_PORT == "" {
		APP_PORT = "3000"
	}

	// dbPool, err := initDatabase()
	// if err != nil {
	// 	slog.Error("Database error:", "error", err)
	// 	return
	// }

	// err = runMigrations(dbPool)
	// if err != nil {
	// 	slog.Error("Database migration error:", "error", err)
	// 	return
	// }

	hub := client.NewHub()
	go hub.Run()

	// Command router is used to execute commands from the client
	commandRouter := client.NewCommandRouter()
	registerCommands(commandRouter)

	handlers.RegisterHandlers(router, hub, commandRouter)

	slog.Info("Database initialized")

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", APP_PORT),
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Error starting server:", "error", err)
			return
		}
	}()

	slog.Info(fmt.Sprintf("Server started at %s ", srv.Addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down server:", "error", err)
		return
	}
	slog.Info("Server stopped")
}

func initDatabase() (*gorm.DB, error) {
	// TODO add env variables and take it from there
	// dbConn := database.DatabaseConnector{
	// 	Host:     "localhost",
	// 	Port:     5432,
	// 	Username: "test",
	// 	Password: "test",
	// 	DBName:   "go_todo_lists",
	// }

	// dbPool, err := dbConn.Connect2()
	// if err != nil {
	// 	return nil, err
	// }

	// return dbPool, nil
	return nil, nil
}

func runMigrations(db *gorm.DB) error {
	// err := database.Up(db)
	// if err != nil {
	// 	fmt.Printf("Error running migrations: %v\n", err)
	// 	return err
	// }

	// fmt.Println("Migrations run successfully")
	return nil
}

func registerCommands(commandRouter *client.CommandRouter) {
	err := commandRouter.RegisterCommand(commands.NewTestCommand())
	err = commandRouter.RegisterCommand(commands.NewTokenAddedCommand())
	err = commandRouter.RegisterCommand(commands.NewTokenMovedCommand())
	err = commandRouter.RegisterCommand(commands.NewMapChangedCommand())
	if err != nil {
		slog.Error("Error registering commands:", "error", err)
		panic(err)
	}
}
