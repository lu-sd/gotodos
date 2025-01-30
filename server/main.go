package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
)

type Todo struct {
	ID   string `json:"id"`
	Done bool   `json:"done"`
	Body string `json:"body"`
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	todos := []Todo{}
	// global State (todos list might not persist across multiple instances)
	// If todos is a global slice, it will reset when the app restarts.
	// Fix: Store the todos in a database (like SQLite, PostgreSQL, or Redis)

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.JSON(todos)
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		// Why Use & (Pointer)?
		// Efficient Memory Usage
		// Instead of copying the entire struct, using a pointer means you’re passing around a reference to the data, which is more efficient.
		//
		//Modifying the Original Struct
		// If BodyParser(todo) modifies todo, it directly modifies the struct in memory rather than a copy.
		if err := c.BodyParser(todo); err != nil {
			// The todo argument must be a pointer (&Todo{}) so that BodyParser can modify the struct directly.
			return err
			// parses the request body into a Go struct, return error if the request body can't be parsed
		}
		todo.ID = uuid.New().String()
		todos = append(todos, *todo)
		return c.JSON(todos)
	})
	app.Patch("/api/todos/:id/done", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, t := range todos {
			if t.ID == id {
				todos[i].Done = true
				break
			}
		}
		return c.JSON(todos)
	})
	app.Listen(":3000")
}
