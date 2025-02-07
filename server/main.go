package main

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/lu-sd/gotodos/query"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect to SQLite database
	db, err := sql.Open("sqlite3", "db/database.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize sqlc queries
	queries := query.New(db)

	app := fiber.New()
	distPath := "../client/dist/"
	app.Static("/", distPath)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Get all todos
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		todos, err := queries.GetTodos(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch todos"})
		}
		return c.JSON(todos)
	})

	// Create a new todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		var todo Todo
		if err := c.BodyParser(&todo); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Insert into DB
		res, err := queries.CreateTodo(c.Context(), query.CreateTodoParams{
			Title: todo.Title,
			Done:  false,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to insert todo"})
		}

		// Retrieve ID from result
		id, _ := res.LastInsertId()
		todo.ID = int(id)
		todo.Done = false

		return c.JSON(todo)
	})

	// Toggle todo done status
	app.Patch("/api/todos/:id/done", func(c *fiber.Ctx) error {
		idstr := c.Params("id")
		// Convert id to int64
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}
		// Get the current todo
		todo, err := queries.GetTodoByID(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
		}

		// Toggle state
		newState := !todo.Done
		err = queries.UpdateTodoStatus(c.Context(), query.UpdateTodoStatusParams{
			Done: newState,
			ID:   todo.ID,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update todo"})
		}

		// Return updated todo
		todo.Done = newState
		return c.JSON(todo)
	})

	// Delete a todo
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		idstr := c.Params("id")
		// Convert id to int64
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		// Get the todo before deleting
		todo, err := queries.GetTodoByID(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
		}

		// Delete the todo
		err = queries.DeleteTodo(c.Context(), id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete todo"})
		}

		return c.JSON(todo)
	})

	log.Fatal(app.Listen(":3000"))
}
