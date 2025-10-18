package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jphacks/os_2502/back/api/internal/worker"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 本番環境では適切なオリジンチェックを実装
		return true
	},
}

// WebSocketHandler WebSocketハンドラー
type WebSocketHandler struct {
	monitor *worker.UploadMonitor
	clients map[string]map[*websocket.Conn]bool
	mu      sync.RWMutex
}

// NewWebSocketHandler WebSocketハンドラーを作成
func NewWebSocketHandler(monitor *worker.UploadMonitor) *WebSocketHandler {
	return &WebSocketHandler{
		monitor: monitor,
		clients: make(map[string]map[*websocket.Conn]bool),
	}
}

// HandleUploadStatus アップロード状況のWebSocket接続
func (h *WebSocketHandler) HandleUploadStatus(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		http.Error(w, "group_id is required", http.StatusBadRequest)
		return
	}

	// WebSocketにアップグレード
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// クライアントを登録
	h.registerClient(groupID, conn)
	defer h.unregisterClient(groupID, conn)

	log.Printf("WebSocket client connected for group %s", groupID)

	// コンテキストを作成
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// ハートビート用のゴルーチン
	go h.heartbeat(ctx, conn)

	// 定期的にステータスを送信
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// ステータスを取得
			status, err := h.monitor.CheckUploadStatus(ctx, groupID)
			if err != nil {
				log.Printf("Failed to check upload status: %v", err)
				continue
			}

			// JSON に変換
			data, err := json.Marshal(status)
			if err != nil {
				log.Printf("Failed to marshal status: %v", err)
				continue
			}

			// クライアントに送信
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Failed to send message: %v", err)
				return
			}

			// 全員アップロード完了の場合は接続を閉じる
			if status.Uploaded >= status.Total && status.Status == "completed" {
				log.Printf("All uploads completed for group %s, closing connection", groupID)

				// 完了メッセージを送信
				completedMsg := map[string]interface{}{
					"type":    "completed",
					"message": "コラージュ生成中...",
				}
				msgData, _ := json.Marshal(completedMsg)
				conn.WriteMessage(websocket.TextMessage, msgData)

				time.Sleep(1 * time.Second)
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// registerClient クライアントを登録
func (h *WebSocketHandler) registerClient(groupID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[groupID] == nil {
		h.clients[groupID] = make(map[*websocket.Conn]bool)
	}
	h.clients[groupID][conn] = true
}

// unregisterClient クライアントを登録解除
func (h *WebSocketHandler) unregisterClient(groupID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[groupID] != nil {
		delete(h.clients[groupID], conn)
		if len(h.clients[groupID]) == 0 {
			delete(h.clients, groupID)
		}
	}
}

// BroadcastToGroup グループの全クライアントにメッセージを送信
func (h *WebSocketHandler) BroadcastToGroup(groupID string, message interface{}) {
	h.mu.RLock()
	clients := h.clients[groupID]
	h.mu.RUnlock()

	if clients == nil {
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Failed to send message to client: %v", err)
		}
	}
}

// heartbeat ハートビート
func (h *WebSocketHandler) heartbeat(ctx context.Context, conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("Failed to send ping: %v", err)
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// HandleStatus ステータス確認用のHTTPエンドポイント（ポーリング用）
func (h *WebSocketHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		respondError(w, http.StatusBadRequest, "group_idが必要です")
		return
	}

	// ステータスを取得
	status, err := h.monitor.CheckUploadStatus(r.Context(), groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "ステータスの取得に失敗しました")
		return
	}

	respondJSON(w, http.StatusOK, status)
}
