package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/alexsasharegan/dotenv"
	"github.com/heraju/mestri"
	"github.com/heraju/mestri/app"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	db, err := sql.Open("postgres", mestri.PsqlInfo)
	die(err)
	defer db.Close()
	// Echo instance
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	app.Connect(e, db)
	// Start server
	port := os.Getenv("PORT")

	if port == "" {
		//log.Fatal("$PORT must be set")
		err := dotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		port = os.Getenv("PORT")
	}
	e.Logger.Fatal(e.Start(":" + port))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// Utility functions
func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
