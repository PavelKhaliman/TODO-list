package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

// InitDB - подключение к БД
func InitDB() error {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {

		databaseUrl = "postgres://pavel:123@localhost:5432/todo"
	}
	var err error
	pool, err = pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		return err
	}
	return createTable()
}

// CloseDB - закрывает соединение с базой
func CloseDB() {
	if pool != nil {
		pool.Close()
	}
}

// createTable - создание таблицы если не существует
func createTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        description TEXT,
        status TEXT CHECK (status IN ('new','in_progress','done')) DEFAULT 'new',
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );
    `
	_, err := pool.Exec(context.Background(), query)
	return err
}

// GetPool возвращает пул соединений
func GetPool() *pgxpool.Pool {
	return pool
}
