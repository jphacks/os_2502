package internal

import (
	"database/sql"
	"log"

	"github.com/noonyuu/collage/api/internal/db"
)

// Server represents the main application server
type Server struct {
	DB *sql.DB
}

// NewServer creates a new server instance with dependencies
func NewServer() (*Server, error) {
	// Initialize database connection
	dbConfig := db.MySQLConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "collage",
		Username: "root",
		Password: "password",
	}

	database, err := db.NewMySQLConnection(dbConfig)
	if err != nil {
		return nil, err
	}

	return &Server{
		DB: database,
	}, nil
}

// Start starts the server and its dependencies
func (s *Server) Start() error {
	log.Println("Server started successfully")
	return nil
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	// Close database connection
	if err := s.DB.Close(); err != nil {
		return err
	}

	log.Println("Server stopped successfully")
	return nil
}
