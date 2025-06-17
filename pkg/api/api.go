package api

import (
	"strconv"
	"time"
	"todo/pkg/db"

	"github.com/gofiber/fiber/v2"
)

// Task - структура задачи
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RegisterRoutes - регистрация endpoints для API
func RegisterRoutes(app *fiber.App) {
	app.Post("/tasks", createTask)
	app.Get("/tasks", getTasks)
	app.Put("/tasks/:id", updateTask)
	app.Delete("/tasks/:id", deleteTask)
}

// createTask - создание новой задачи
func createTask(c *fiber.Ctx) error {
	var task Task
	if err := c.BodyParser(&task); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ввод"})
	}

	query := `
    INSERT INTO tasks (title, description, status, created_at, updated_at)
    VALUES ($1, $2, $3, NOW(), NOW())
    RETURNING id, title, description, status, created_at, updated_at
    `
	err := db.GetPool().QueryRow(c.Context(), query, task.Title, task.Description, task.Status).Scan(
		&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка создания задачи"})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

// getTasks - возвращает список задач
func getTasks(c *fiber.Ctx) error {
	rows, err := db.GetPool().Query(c.Context(), "SELECT id, title, description, status, created_at, updated_at FROM tasks")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка получения задач"})
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка при парсинге задач"})
		}
		tasks = append(tasks, t)
	}
	return c.JSON(tasks)
}

// updateTask обновляет задачу по id
func updateTask(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ошибка ID задачи"})
	}

	var task Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный ввод"})
	}

	query := `
    UPDATE tasks
    SET title = $1, description = $2, status = $3, updated_at = NOW()
    WHERE id = $4
    RETURNING created_at, updated_at
    `
	err = db.GetPool().QueryRow(c.Context(), query, task.Title, task.Description, task.Status, id).Scan(&task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обновления задачи"})
	}
	task.ID = id
	return c.JSON(task)
}

// deleteTask удаляет задачу по id
func deleteTask(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ошибка ID задачи"})
	}

	result, err := db.GetPool().Exec(c.Context(), "DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка удаления задачи"})
	}
	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Задача не найдена"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
