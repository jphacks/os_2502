package internal

import (
	"database/sql"
	"net/http"

	"github.com/jphacks/os_2502/back/api/internal/handler"
	"github.com/jphacks/os_2502/back/api/internal/infrastructure/repository"
	"github.com/jphacks/os_2502/back/api/internal/usecase"
	"github.com/jphacks/os_2502/back/api/internal/worker"
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

	// Repository 初期化
	userRepo := repository.NewUserRepository(r.db)
	groupRepo := repository.NewGroupRepositorySQLBoiler(r.db)
	groupMemberRepo := repository.NewGroupMemberRepositorySQLBoiler(r.db)
	friendRepo := repository.NewFriendRepositorySQLBoiler(r.db)
	deviceTokenRepo := repository.NewDeviceTokenRepositorySQLBoiler(r.db)
	collageTemplateRepo := repository.NewCollageTemplateRepositorySQLBoiler(r.db)
	collageResultRepo := repository.NewCollageResultRepositorySQLBoiler(r.db)
	uploadImageRepo := repository.NewUploadImageRepositorySQLBoiler(r.db)
	resultDownloadRepo := repository.NewResultDownloadRepositorySQLBoiler(r.db)
	templatePartRepo := repository.NewTemplatePartRepository(r.db)
	groupPartAssignmentRepo := repository.NewGroupPartAssignmentRepository(r.db)
	uploadImagesCollageResultRepo := repository.NewUploadImagesCollageResultRepository(r.db)

	// UseCase 初期化
	userUC := usecase.NewUserUseCase(userRepo)
	groupUC := usecase.NewGroupUseCase(groupRepo, groupMemberRepo)
	friendUC := usecase.NewFriendUseCase(friendRepo)
	deviceTokenUC := usecase.NewDeviceTokenUseCase(deviceTokenRepo)
	collageTemplateUC := usecase.NewCollageTemplateUseCase(collageTemplateRepo)
	collageResultUC := usecase.NewCollageResultUseCase(collageResultRepo)
	uploadImageUC := usecase.NewUploadImageUseCase(uploadImageRepo)
	resultDownloadUC := usecase.NewResultDownloadUseCase(resultDownloadRepo)
	templatePartUC := usecase.NewTemplatePartUseCase(templatePartRepo)
	groupPartAssignmentUC := usecase.NewGroupPartAssignmentUseCase(groupPartAssignmentRepo)
	uploadImagesCollageResultUC := usecase.NewUploadImagesCollageResultUseCase(uploadImagesCollageResultRepo)

	// Worker 初期化
	uploadMonitor := worker.NewUploadMonitor(uploadImageRepo)

	// Handler 初期化
	userHandler := handler.NewUserHandler(userUC)
	groupHandler := handler.NewGroupHandler(groupUC)
	friendHandler := handler.NewFriendHandler(friendUC)
	deviceTokenHandler := handler.NewDeviceTokenHandler(deviceTokenUC)
	collageTemplateHandler := handler.NewCollageTemplateHandler(collageTemplateUC)
	collageResultHandler := handler.NewCollageResultHandler(collageResultUC)
	uploadImageHandler := handler.NewUploadImageHandler(uploadImageUC)
	resultDownloadHandler := handler.NewResultDownloadHandler(resultDownloadUC)
	templatePartHandler := handler.NewTemplatePartHandler(templatePartUC)
	groupPartAssignmentHandler := handler.NewGroupPartAssignmentHandler(groupPartAssignmentUC)
	uploadImagesCollageResultHandler := handler.NewUploadImagesCollageResultHandler(uploadImagesCollageResultUC)
	websocketHandler := handler.NewWebSocketHandler(uploadMonitor)

	// User エンドポイント
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			userHandler.CreateUser(w, r)
		case http.MethodGet:
			userHandler.ListUsers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/users/firebase", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.GetUserByFirebaseUID(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUser(w, r)
		case http.MethodPut, http.MethodPatch:
			userHandler.UpdateUser(w, r)
		case http.MethodDelete:
			userHandler.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Group エンドポイント
	mux.HandleFunc("/api/groups", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			groupHandler.CreateGroup(w, r)
		case http.MethodGet:
			groupHandler.ListGroups(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/groups/", groupHandler.GetGroupByID)

	// Friend エンドポイント
	mux.HandleFunc("/api/friends", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			friendHandler.SendFriendRequest(w, r)
		case http.MethodGet:
			friendHandler.GetFriends(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Device Token エンドポイント
	mux.HandleFunc("/api/device-tokens", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			deviceTokenHandler.RegisterDeviceToken(w, r)
		case http.MethodGet:
			deviceTokenHandler.GetUserDeviceTokens(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Collage Template エンドポイント
	mux.HandleFunc("/api/templates", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			collageTemplateHandler.CreateTemplate(w, r)
		case http.MethodGet:
			collageTemplateHandler.ListTemplates(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/templates/", collageTemplateHandler.GetTemplate)

	// Collage Result エンドポイント
	mux.HandleFunc("/api/results", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			collageResultHandler.CreateResult(w, r)
		case http.MethodGet:
			collageResultHandler.GetResultsByGroup(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/results/", collageResultHandler.GetResult)

	// Upload Image エンドポイント
	mux.HandleFunc("/api/images", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			uploadImageHandler.UploadImage(w, r)
		case http.MethodGet:
			uploadImageHandler.GetImagesByGroup(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/images/", uploadImageHandler.GetImage)

	// Result Download エンドポイント
	mux.HandleFunc("/api/downloads", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			resultDownloadHandler.RecordDownload(w, r)
		case http.MethodGet:
			resultDownloadHandler.GetDownloadsByResult(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Template Part エンドポイント
	mux.HandleFunc("/api/template-parts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			templatePartHandler.CreateTemplatePart(w, r)
		case http.MethodGet:
			templatePartHandler.ListTemplateParts(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/template-parts/", templatePartHandler.GetTemplatePart)

	// Group Part Assignment エンドポイント
	mux.HandleFunc("/api/part-assignments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			groupPartAssignmentHandler.CreateGroupPartAssignment(w, r)
		case http.MethodGet:
			groupPartAssignmentHandler.ListGroupPartAssignments(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/part-assignments/", groupPartAssignmentHandler.GetGroupPartAssignment)

	// Upload Images Collage Result エンドポイント
	mux.HandleFunc("/api/image-results", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			uploadImagesCollageResultHandler.CreateUploadImagesCollageResult(w, r)
		case http.MethodGet:
			uploadImagesCollageResultHandler.ListUploadImagesCollageResults(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// WebSocket エンドポイント
	mux.HandleFunc("/api/ws/upload-status", websocketHandler.HandleUploadStatus)
	mux.HandleFunc("/api/status", websocketHandler.HandleStatus)

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
  
	return middleware.CORSMiddleware(mux)
}
