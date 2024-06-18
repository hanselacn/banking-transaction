package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	cfg "github.com/hanselacn/banking-transaction/config"
	"github.com/hanselacn/banking-transaction/internal/entity"
	"github.com/hanselacn/banking-transaction/internal/handler"
	"github.com/hanselacn/banking-transaction/internal/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestInvalidUserName(t *testing.T) {
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
		Worker: cfg.Worker{
			PayoutInterval: os.Getenv("PAYOUT_INTERVAL"),
			PayoutTimeUnit: os.Getenv("PAYOUT_TIME_UNIT"),
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

	r.Handle("/ping", http.HandlerFunc(pingHandler))
	handler.MountUserHandler(r, h, m)
	handler.MountAccountHandler(r, h, m)

	reqBody := entity.CreateUserInput{
		Username: "</script>",
		Fullname: "yayayay",
		Password: "okeL@hKalauB3gitu",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UsersHandler.CreateUserHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.Equal(t, "success create user", respBody["errors"])
	assert.NoError(t, err)
}

func TestInvalidPassword(t *testing.T) {
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
		Worker: cfg.Worker{
			PayoutInterval: os.Getenv("PAYOUT_INTERVAL"),
			PayoutTimeUnit: os.Getenv("PAYOUT_TIME_UNIT"),
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

	r.Handle("/ping", http.HandlerFunc(pingHandler))
	handler.MountUserHandler(r, h, m)
	handler.MountAccountHandler(r, h, m)

	reqBody := entity.CreateUserInput{
		Username: "IniAdalahPercobaan",
		Fullname: "yayayay",
		Password: "1.1.1",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UsersHandler.CreateUserHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.Equal(t, "success create user", respBody["errors"])
	assert.NoError(t, err)
}

func TestInvalidAuthorization(t *testing.T) {
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
		Worker: cfg.Worker{
			PayoutInterval: os.Getenv("PAYOUT_INTERVAL"),
			PayoutTimeUnit: os.Getenv("PAYOUT_TIME_UNIT"),
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

	r.Handle("/ping", http.HandlerFunc(pingHandler))
	handler.MountUserHandler(r, h, m)
	handler.MountAccountHandler(r, h, m)

	reqBody := entity.CreateUserInput{
		Username: "</script>",
		Fullname: "yayayay",
		Password: "[]121",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := m.AuthenticationMiddleware(http.HandlerFunc(h.UsersHandler.CreateUserHandler))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.Equal(t, "success create user", respBody["message"])
	assert.NoError(t, err)
}

func TestSQLInjection(t *testing.T) {
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
		Worker: cfg.Worker{
			PayoutInterval: os.Getenv("PAYOUT_INTERVAL"),
			PayoutTimeUnit: os.Getenv("PAYOUT_TIME_UNIT"),
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

	r.Handle("/ping", http.HandlerFunc(pingHandler))
	handler.MountUserHandler(r, h, m)
	handler.MountAccountHandler(r, h, m)

	reqBody := entity.User{
		Username: "a OR 1=1",
		Role:     "super_admin",
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("failed marshall body")
	}

	req := httptest.NewRequest("PUT", "/", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UsersHandler.UpdateRoleByUserName)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var respBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)
	assert.Equal(t, "success update user", respBody["errors"])
	assert.NoError(t, err)
}
