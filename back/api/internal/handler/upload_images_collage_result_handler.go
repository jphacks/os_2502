package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/upload_images_collage_result"
	"github.com/jphacks/os_2502/back/api/internal/usecase"
)

type UploadImagesCollageResultHandler struct {
	useCase *usecase.UploadImagesCollageResultUseCase
}

func NewUploadImagesCollageResultHandler(useCase *usecase.UploadImagesCollageResultUseCase) *UploadImagesCollageResultHandler {
	return &UploadImagesCollageResultHandler{useCase: useCase}
}

type CreateUploadImagesCollageResultRequest struct {
	ImageID   string `json:"image_id"`
	ResultID  string `json:"result_id"`
	PositionX int    `json:"position_x"`
	PositionY int    `json:"position_y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	SortOrder int    `json:"sort_order"`
}

type UpdateUploadImagesCollageResultPositionRequest struct {
	PositionX int `json:"position_x"`
	PositionY int `json:"position_y"`
	Width     int `json:"width"`
	Height    int `json:"height"`
}

type UpdateUploadImagesCollageResultSortOrderRequest struct {
	SortOrder int `json:"sort_order"`
}

type UploadImagesCollageResultResponse struct {
	ImageID   string `json:"image_id"`
	ResultID  string `json:"result_id"`
	PositionX int    `json:"position_x"`
	PositionY int    `json:"position_y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	SortOrder int    `json:"sort_order"`
	CreatedAt string `json:"created_at"`
}

// UploadImagesCollageResultエンティティをUploadImagesCollageResultResponseに変換
func toUploadImagesCollageResultResponse(uicr *upload_images_collage_result.UploadImagesCollageResult) UploadImagesCollageResultResponse {
	return UploadImagesCollageResultResponse{
		ImageID:   uicr.ImageID().String(),
		ResultID:  uicr.ResultID().String(),
		PositionX: uicr.PositionX(),
		PositionY: uicr.PositionY(),
		Width:     uicr.Width(),
		Height:    uicr.Height(),
		SortOrder: uicr.SortOrder(),
		CreatedAt: uicr.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *UploadImagesCollageResultHandler) CreateUploadImagesCollageResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	var req CreateUploadImagesCollageResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	imageID, err := uuid.Parse(req.ImageID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な画像IDです")
		return
	}

	resultID, err := uuid.Parse(req.ResultID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な結果IDです")
		return
	}

	uicr, err := h.useCase.CreateUploadImagesCollageResult(
		r.Context(),
		imageID,
		resultID,
		req.PositionX,
		req.PositionY,
		req.Width,
		req.Height,
		req.SortOrder,
	)
	if err != nil {
		switch err {
		case upload_images_collage_result.ErrInvalidImageID, upload_images_collage_result.ErrInvalidResultID, upload_images_collage_result.ErrInvalidDimensions:
			respondError(w, http.StatusBadRequest, err.Error())
		case upload_images_collage_result.ErrUploadImagesCollageResultAlreadyExists:
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "画像コラージュ結果の作成に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusCreated, toUploadImagesCollageResultResponse(uicr))
}

func (h *UploadImagesCollageResultHandler) GetUploadImagesCollageResultByImageIDAndResultID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	imageIDStr := r.URL.Query().Get("image_id")
	if imageIDStr == "" {
		respondError(w, http.StatusBadRequest, "画像IDが必要です")
		return
	}

	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な画像IDです")
		return
	}

	resultIDStr := r.URL.Query().Get("result_id")
	if resultIDStr == "" {
		respondError(w, http.StatusBadRequest, "結果IDが必要です")
		return
	}

	resultID, err := uuid.Parse(resultIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な結果IDです")
		return
	}

	uicr, err := h.useCase.GetUploadImagesCollageResultByImageIDAndResultID(r.Context(), imageID, resultID)
	if err != nil {
		if err == upload_images_collage_result.ErrUploadImagesCollageResultNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "画像コラージュ結果の取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toUploadImagesCollageResultResponse(uicr))
}

func (h *UploadImagesCollageResultHandler) UpdateUploadImagesCollageResultPosition(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	imageIDStr := r.URL.Query().Get("image_id")
	if imageIDStr == "" {
		respondError(w, http.StatusBadRequest, "画像IDが必要です")
		return
	}

	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な画像IDです")
		return
	}

	resultIDStr := r.URL.Query().Get("result_id")
	if resultIDStr == "" {
		respondError(w, http.StatusBadRequest, "結果IDが必要です")
		return
	}

	resultID, err := uuid.Parse(resultIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な結果IDです")
		return
	}

	var req UpdateUploadImagesCollageResultPositionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	uicr, err := h.useCase.UpdateUploadImagesCollageResultPosition(r.Context(), imageID, resultID, req.PositionX, req.PositionY, req.Width, req.Height)
	if err != nil {
		switch err {
		case upload_images_collage_result.ErrInvalidDimensions:
			respondError(w, http.StatusBadRequest, err.Error())
		case upload_images_collage_result.ErrUploadImagesCollageResultNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "画像コラージュ結果の位置更新に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toUploadImagesCollageResultResponse(uicr))
}

func (h *UploadImagesCollageResultHandler) UpdateUploadImagesCollageResultSortOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	imageIDStr := r.URL.Query().Get("image_id")
	if imageIDStr == "" {
		respondError(w, http.StatusBadRequest, "画像IDが必要です")
		return
	}

	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な画像IDです")
		return
	}

	resultIDStr := r.URL.Query().Get("result_id")
	if resultIDStr == "" {
		respondError(w, http.StatusBadRequest, "結果IDが必要です")
		return
	}

	resultID, err := uuid.Parse(resultIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な結果IDです")
		return
	}

	var req UpdateUploadImagesCollageResultSortOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	uicr, err := h.useCase.UpdateUploadImagesCollageResultSortOrder(r.Context(), imageID, resultID, req.SortOrder)
	if err != nil {
		if err == upload_images_collage_result.ErrUploadImagesCollageResultNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "画像コラージュ結果のソート順更新に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toUploadImagesCollageResultResponse(uicr))
}

func (h *UploadImagesCollageResultHandler) DeleteUploadImagesCollageResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	imageIDStr := r.URL.Query().Get("image_id")
	if imageIDStr == "" {
		respondError(w, http.StatusBadRequest, "画像IDが必要です")
		return
	}

	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な画像IDです")
		return
	}

	resultIDStr := r.URL.Query().Get("result_id")
	if resultIDStr == "" {
		respondError(w, http.StatusBadRequest, "結果IDが必要です")
		return
	}

	resultID, err := uuid.Parse(resultIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な結果IDです")
		return
	}

	if err := h.useCase.DeleteUploadImagesCollageResult(r.Context(), imageID, resultID); err != nil {
		if err == upload_images_collage_result.ErrUploadImagesCollageResultNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "画像コラージュ結果の削除に失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UploadImagesCollageResultHandler) ListUploadImagesCollageResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// クエリパラメータからlimitとoffsetを取得
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	results, err := h.useCase.ListUploadImagesCollageResults(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "画像コラージュ結果一覧の取得に失敗しました")
		return
	}

	responses := make([]UploadImagesCollageResultResponse, 0, len(results))
	for _, uicr := range results {
		responses = append(responses, toUploadImagesCollageResultResponse(uicr))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"upload_images_collage_results": responses,
		"limit":                         limit,
		"offset":                        offset,
		"count":                         len(responses),
	})
}

func (h *UploadImagesCollageResultHandler) GetUploadImagesCollageResultsByImageID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	imageIDStr := r.URL.Query().Get("image_id")
	if imageIDStr == "" {
		respondError(w, http.StatusBadRequest, "画像IDが必要です")
		return
	}

	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な画像IDです")
		return
	}

	results, err := h.useCase.GetUploadImagesCollageResultsByImageID(r.Context(), imageID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "画像コラージュ結果の取得に失敗しました")
		return
	}

	responses := make([]UploadImagesCollageResultResponse, 0, len(results))
	for _, uicr := range results {
		responses = append(responses, toUploadImagesCollageResultResponse(uicr))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"upload_images_collage_results": responses,
		"count":                         len(responses),
	})
}

func (h *UploadImagesCollageResultHandler) GetUploadImagesCollageResultsByResultID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	resultIDStr := r.URL.Query().Get("result_id")
	if resultIDStr == "" {
		respondError(w, http.StatusBadRequest, "結果IDが必要です")
		return
	}

	resultID, err := uuid.Parse(resultIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な結果IDです")
		return
	}

	results, err := h.useCase.GetUploadImagesCollageResultsByResultID(r.Context(), resultID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "画像コラージュ結果の取得に失敗しました")
		return
	}

	responses := make([]UploadImagesCollageResultResponse, 0, len(results))
	for _, uicr := range results {
		responses = append(responses, toUploadImagesCollageResultResponse(uicr))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"upload_images_collage_results": responses,
		"count":                         len(responses),
	})
}
