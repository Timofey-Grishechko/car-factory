package main

import (
	"car_factory/internal/db"
	"car_factory/internal/server"
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	ctx := context.Background()
	pool, err := db.CreateConnection(ctx)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных")
	}
	fmt.Println("Успешное подключение к Postgre")
	defer pool.Close()

	server.StartServer(pool)
}
