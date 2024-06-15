package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"alpha.com/configuration"
	"alpha.com/internal/alpha.com/pkg/mongodb"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	app *fiber.App
}

func NewServer(app *fiber.App) *server {
	return &server{
		app: app,
	}
}

func (s *server) StartHttpServer(mongoClient *mongo.Client) {
	go func() {
		gracefulShutdown(s.app, mongoClient)
	}()

	if err := s.app.Listen(fmt.Sprintf(":%s", configuration.Port)); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Cannot start server - ERROR: %v\n", err)
		panic("cannot start server")
	}
}

func gracefulShutdown(app *fiber.App, mongoClient *mongo.Client) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutdown Server")

	_, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		fmt.Printf("Server Shutdown Error: %v\n", err)
	}

	fmt.Println("Fiber server gracefully stopped.")

	mongodb.DisconnectMongoDB(mongoClient)

	fmt.Println("Server exiting")
}
