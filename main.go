package main

import (
	"fmt"

	"todo/pkg/db"
	"todo/pkg/server"
)

func main() {
	// Инициализация БД
	if err := db.InitDB(); err != nil {
		fmt.Println("ошибка при инициализации БД", err)
	}
	defer db.CloseDB()

	// Запуск сервера
	s := server.NewServer()
	if err := s.Listen(":3000"); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
