package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jphacks/os_2502/back/api/config"
	"github.com/jphacks/os_2502/back/api/internal"
	"github.com/jphacks/os_2502/back/api/internal/db"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/repository"
	"github.com/jphacks/os_2502/back/api/internal/worker"
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

	// コラージュ生成ワーカーを起動
	groupRepo := repository.NewGroupRepositorySQLBoiler(database)
	groupMemberRepo := repository.NewGroupMemberRepositorySQLBoiler(database)
	collageGenerator := worker.NewCollageGenerator(groupRepo, groupMemberRepo, 10*time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go collageGenerator.Start(ctx)

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
	cancel() // ワーカーを停止

	log.Println("シャットダウン完了")
}
