package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/noonyuu/collage/api/config"
	"github.com/noonyuu/collage/api/internal"
	"github.com/noonyuu/collage/api/internal/db"
)

func main() {
	cfg := config.Load()

	// DB設定の読み込み
	dbConfig := db.MySQLConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Database: cfg.Database.Name,
		Username: cfg.Database.User,
		Password: cfg.Database.Password,
	}

	database, err := db.NewMySQLConnection(dbConfig)
	if err != nil {
		log.Fatalf("データベース接続に失敗: %v", err)
	}
	defer database.Close()

	// ルーターの初期化と設定
	router := internal.NewRouter(database)
	handler := router.SetupRoutes()

	// サーバーを起動
	go func() {
		log.Println("Starting server on :8080")
		if err := http.ListenAndServe(":8080", handler); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started successfully")
	log.Println("- API: http://localhost:8080")

	// グレースフルシャットダウン
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("サーバーをシャットダウン中...")

	log.Println("シャットダウン完了")
}
