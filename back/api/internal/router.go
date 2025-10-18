package internal

import (
	"database/sql"
	"net/http"
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

	return mux
}
