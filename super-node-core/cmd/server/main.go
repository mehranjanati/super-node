package main

import (
	"log"
	"super-node-core/internal/cocoon"
	"super-node-core/internal/gateway"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Initialize Monolithic Core Components
	log.Println("[SuperNode] Initializing Cocoon Runtime...")
	orchestrator := cocoon.NewOrchestrator()

	log.Println("[SuperNode] Initializing Gateway...")
	handler := gateway.NewHandler(orchestrator)

	// Setup Web Server (Fiber)
	app := fiber.New()

	// Enable CORS for SvelteKit
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, http://localhost:5174",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Routes
	api := app.Group("/api")
	api.Post("/transaction", handler.SubmitTransaction)

	// Start Server
	log.Println("[SuperNode] Ready to process transactions on :3001")
	log.Fatal(app.Listen(":3001"))
}
