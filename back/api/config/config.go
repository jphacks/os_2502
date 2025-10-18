package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

type ServerConfig struct {
	Port int
}

func Load() *Config {
	// .envファイルから環境変数を読み込み
	loadEnvFile()

	mysqlHost := getEnvOrDefault("MYSQL_HOST", "mysql:3306")
	host, port := parseHostPort(mysqlHost)

	dbPort, _ := strconv.Atoi(getEnvOrDefault("DB_PORT", port))
	serverPort, _ := strconv.Atoi(getEnvOrDefault("SERVER_PORT", "8080"))

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", host),
			Port:     dbPort,
			Name:     getEnvOrDefault("MYSQL_DATABASE", ""),
			User:     getEnvOrDefault("MYSQL_USER", ""),
			Password: getEnvOrDefault("MYSQL_PASSWORD", ""),
		},
		Server: ServerConfig{
			Port: serverPort,
		},
	}
}

// .envファイルを読み込んで環境変数に設定
func loadEnvFile() {
	envFile := ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		log.Printf(".env ファイルが見つかりません: %s", envFile)
		return
	}

	file, err := os.Open(envFile)
	if err != nil {
		log.Printf(".env ファイルの読み込みに失敗: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 空行やコメント行をスキップ
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// KEY=VALUE 形式をパース
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 既に環境変数が設定されている場合は上書きしない
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf(".env ファイルの読み込み中にエラー: %v", err)
	}
}

// ホスト:ポートの形式を分離
func parseHostPort(hostPort string) (host, port string) {
	parts := strings.Split(hostPort, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return hostPort, "3306"
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
