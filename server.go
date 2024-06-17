package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	cfg "github.com/hanselacn/banking-transaction/config"
	"github.com/hanselacn/banking-transaction/internal/handler"
	"github.com/hanselacn/banking-transaction/internal/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Create a new ServeMux
	var (
		eventName = "server.start"
	)
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(eventName, err)
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	cfg := cfg.Config{
		DB: cfg.Database{
			Driver:   os.Getenv("DB_DRIVER"),
			Name:     os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
		},
		Server: cfg.Server{
			Port: os.Getenv("API_PORT"),
			TLS:  os.Getenv("API_TLS"),
		},
	}

	connStr := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=disable", cfg.DB.Driver, cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Name)
	db, err := sql.Open(cfg.DB.Driver, connStr)
	if err != nil {
		log.Println(eventName, err)
		fmt.Println("Error starting database connection:", err)
		log.Fatal(err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Println(eventName, err)
		fmt.Println("Error ping database connection:", err)
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database!")

	h := handler.NewHandler(db)
	m := middleware.NewMiddleware(db)
	handler.MountUserHandler(r, h, m)
	handler.MountAccountHandler(r, h, m)
	// Start the server
	fmt.Printf("Starting server on :%s \n", cfg.Server.Port)
	switch os.Getenv("API_TLS") {
	case "true":
		fmt.Println("Starting server using TLS...")
		err = http.ListenAndServeTLS(fmt.Sprintf(":%s", cfg.Server.Port), "cert.pem", "key.pem", r)
		if err != nil {
			log.Println(eventName, err)
			fmt.Println("Error starting server:", err)
		}
	default:
		err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), r)
		if err != nil {
			log.Println(eventName, err)
			fmt.Println("Error starting server:", err)
		}
	}
}
