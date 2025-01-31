package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	InitDatabase()
	defer DB.Close()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	// todos := []Todo{}
	// global State (todos list might not persist across multiple instances)
	// If todos is a global slice, it will reset when the app restarts.
	// Fix: Store the todos in a database (like SQLite, PostgreSQL, or Redis)

	// app.Get("/api/todos", func(c *fiber.Ctx) error {
	// return c.JSON(todos)
	// })
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		rows, err := DB.Query("SELECT id, title, done FROM todos")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch todos"})
		}
		defer rows.Close()

		var todos []Todo

		for rows.Next() {
			var todo Todo
			if err := rows.Scan(&todo.ID, &todo.Title, &todo.Done); err != nil {
				return err
			}
			todos = append(todos, todo)
		}

		return c.JSON(todos)
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		// Modifying the Original Struct
		// If BodyParser(todo) modifies todo, it directly modifies the struct in memory rather than a copy.
		if err := c.BodyParser(todo); err != nil {
			// The todo argument must be a pointer (&Todo{}) so that BodyParser can modify the struct directly.
			return err
			// parses the request body into a Go struct, return error if the request body can't be parsed
		}
		// Insert into SQLite
		res, err := DB.Exec("INSERT INTO todos (title, done) VALUES (?, ?)", todo.Title, false)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to insert todo"})
		}

		id, _ := res.LastInsertId()
		todo.ID = int(id)
		todo.Done = false

		return c.JSON(todo)
	})

	app.Patch("/api/todos/:id/done", func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Get current state
		var done bool
		err := DB.QueryRow("SELECT done FROM todos WHERE id = ?", id).Scan(&done)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
		}

		// Toggle state
		newState := !done
		_, err = DB.Exec("UPDATE todos SET done = ? WHERE id = ?", newState, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update todo"})
		}
		// Fetch the updated todo
		var updatedTodo struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
			Done  bool   `json:"done"`
		}

		err = DB.QueryRow("SELECT id, title, done FROM todos WHERE id = ?", id).Scan(&updatedTodo.ID, &updatedTodo.Title, &updatedTodo.Done)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve updated todo"})
		}

		return c.JSON(updatedTodo)
	})

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		var deletedTodo Todo
		err := DB.QueryRow("SELECT id, title, done FROM todos WHERE id = ?", id).Scan(&deletedTodo.ID, &deletedTodo.Title, &deletedTodo.Done)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
		}
		_, err = DB.Exec("DELETE FROM todos WHERE id = ?", id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete todo"})
		}
		return c.JSON(deletedTodo)
	})

	log.Fatal(app.Listen(":3000"))
}
