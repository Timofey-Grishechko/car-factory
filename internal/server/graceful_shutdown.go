package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func gracefulShutdown(quit chan os.Signal, server *http.Server) {
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Выключение сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Сервер принудительно выключен", err)
	}
	log.Println("Сервер выключен")
}
