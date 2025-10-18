package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
	"github.com/noonyuu/collage/api/internal/domain/device_token"
	"github.com/noonyuu/collage/api/internal/infrastructure/db"
	"github.com/noonyuu/collage/api/internal/infrastructure/models"
)

type DeviceTokenRepositorySQLBoiler struct {
	db *sql.DB
}

func NewDeviceTokenRepositorySQLBoiler(db *sql.DB) device_token.Repository {
	return &DeviceTokenRepositorySQLBoiler{db: db}
}

// Model to Entity conversion
func toDeviceTokenEntity(m *models.DeviceToken) (*device_token.DeviceToken, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(m.UserID)
	if err != nil {
		return nil, err
	}

	var deviceName *string
	if m.DeviceName.Valid {
		deviceName = &m.DeviceName.String
	}

	var lastUsedAt *time.Time
	if m.LastUsedAt.Valid {
		lastUsedAt = &m.LastUsedAt.Time
	}

	return device_token.Reconstruct(
		id,
		userID,
		m.DeviceToken,
		device_token.DeviceType(m.DeviceType),
		deviceName,
		m.IsActive,
		lastUsedAt,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

// Entity to Model conversion
func toDeviceTokenModel(dt *device_token.DeviceToken) *models.DeviceToken {
	model := &models.DeviceToken{
		ID:          dt.ID().String(),
		UserID:      dt.UserID().String(),
		DeviceToken: dt.DeviceToken(),
		DeviceType:  string(dt.DeviceType()),
		IsActive:    dt.IsActive(),
		CreatedAt:   dt.CreatedAt(),
		UpdatedAt:   dt.UpdatedAt(),
	}

	if deviceName := dt.DeviceName(); deviceName != nil {
		model.DeviceName.Valid = true
		model.DeviceName.String = *deviceName
	}

	if lastUsedAt := dt.LastUsedAt(); lastUsedAt != nil {
		model.LastUsedAt.Valid = true
		model.LastUsedAt.Time = *lastUsedAt
	}

	return model
}

func (r *DeviceTokenRepositorySQLBoiler) Create(ctx context.Context, dt *device_token.DeviceToken) error {
	model := toDeviceTokenModel(dt)
	err := model.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		if db.IsDuplicateError(err) {
			return device_token.ErrDeviceTokenAlreadyExists
		}
		return err
	}
	return nil
}

func (r *DeviceTokenRepositorySQLBoiler) FindByID(ctx context.Context, id uuid.UUID) (*device_token.DeviceToken, error) {
	model, err := models.FindDeviceToken(ctx, r.db, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, device_token.ErrDeviceTokenNotFound
		}
		return nil, err
	}
	return toDeviceTokenEntity(model)
}

func (r *DeviceTokenRepositorySQLBoiler) FindByToken(ctx context.Context, token string) (*device_token.DeviceToken, error) {
	model, err := models.DeviceTokens(
		qm.Where("device_token = ?", token),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, device_token.ErrDeviceTokenNotFound
		}
		return nil, err
	}
	return toDeviceTokenEntity(model)
}

func (r *DeviceTokenRepositorySQLBoiler) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*device_token.DeviceToken, error) {
	modelSlice, err := models.DeviceTokens(
		qm.Where("user_id = ?", userID.String()),
		qm.OrderBy("last_used_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	tokens := make([]*device_token.DeviceToken, len(modelSlice))
	for i, model := range modelSlice {
		dt, err := toDeviceTokenEntity(model)
		if err != nil {
			return nil, err
		}
		tokens[i] = dt
	}
	return tokens, nil
}

func (r *DeviceTokenRepositorySQLBoiler) FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*device_token.DeviceToken, error) {
	modelSlice, err := models.DeviceTokens(
		qm.Where("user_id = ? AND is_active = ?", userID.String(), true),
		qm.OrderBy("last_used_at DESC"),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	tokens := make([]*device_token.DeviceToken, len(modelSlice))
	for i, model := range modelSlice {
		dt, err := toDeviceTokenEntity(model)
		if err != nil {
			return nil, err
		}
		tokens[i] = dt
	}
	return tokens, nil
}

func (r *DeviceTokenRepositorySQLBoiler) Update(ctx context.Context, dt *device_token.DeviceToken) error {
	model, err := models.FindDeviceToken(ctx, r.db, dt.ID().String())
	if err != nil {
		if err == sql.ErrNoRows {
			return device_token.ErrDeviceTokenNotFound
		}
		return err
	}

	model.IsActive = dt.IsActive()
	model.UpdatedAt = dt.UpdatedAt()

	if deviceName := dt.DeviceName(); deviceName != nil {
		model.DeviceName.Valid = true
		model.DeviceName.String = *deviceName
	} else {
		model.DeviceName.Valid = false
	}

	if lastUsedAt := dt.LastUsedAt(); lastUsedAt != nil {
		model.LastUsedAt.Valid = true
		model.LastUsedAt.Time = *lastUsedAt
	} else {
		model.LastUsedAt.Valid = false
	}

	_, err = model.Update(ctx, r.db, boil.Whitelist(
		models.DeviceTokenColumns.IsActive,
		models.DeviceTokenColumns.DeviceName,
		models.DeviceTokenColumns.LastUsedAt,
		models.DeviceTokenColumns.UpdatedAt,
	))
	return err
}

func (r *DeviceTokenRepositorySQLBoiler) Delete(ctx context.Context, id uuid.UUID) error {
	model, err := models.FindDeviceToken(ctx, r.db, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return device_token.ErrDeviceTokenNotFound
		}
		return err
	}

	_, err = model.Delete(ctx, r.db)
	return err
}

func (r *DeviceTokenRepositorySQLBoiler) DeactivateOldTokens(ctx context.Context, days int) (int, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	// last_used_at が NULL または cutoff より古いトークンを非アクティブ化
	count, err := models.DeviceTokens(
		qm.Where("is_active = ? AND (last_used_at IS NULL OR last_used_at < ?)", true, cutoffDate),
	).UpdateAll(ctx, r.db, models.M{
		models.DeviceTokenColumns.IsActive:  false,
		models.DeviceTokenColumns.UpdatedAt: time.Now(),
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
