package internal

import (
	"database/sql"
	"net/http"

	"github.com/jphacks/os_2502/back/api/middleware"
)

type Router struct {
	db *sql.DB
}

// 新しいルーターを作成
func NewRouter(db *sql.DB) *Router {
	return &Router{db: db}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	return middleware.CORSMiddleware(mux)
}
