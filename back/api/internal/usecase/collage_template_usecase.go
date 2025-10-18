package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/collage_template"
)

type CollageTemplateUseCase struct {
	repo collage_template.Repository
}

func NewCollageTemplateUseCase(repo collage_template.Repository) *CollageTemplateUseCase {
	return &CollageTemplateUseCase{repo: repo}
}

// CreateTemplate creates a new collage template
func (uc *CollageTemplateUseCase) CreateTemplate(ctx context.Context, name, filePath string) (*collage_template.CollageTemplate, error) {
	// 同名のテンプレートが存在するかチェック
	existing, err := uc.repo.FindByName(ctx, name)
	if err == nil && existing != nil {
		return nil, collage_template.ErrTemplateAlreadyExists
	}

	// 新規作成
	template, err := collage_template.NewCollageTemplate(name, filePath)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, template); err != nil {
		return nil, err
	}

	return template, nil
}

// GetTemplate retrieves a template by ID
func (uc *CollageTemplateUseCase) GetTemplate(ctx context.Context, templateID uuid.UUID) (*collage_template.CollageTemplate, error) {
	return uc.repo.FindByID(ctx, templateID)
}

// ListTemplates retrieves all templates
func (uc *CollageTemplateUseCase) ListTemplates(ctx context.Context, limit, offset int) ([]*collage_template.CollageTemplate, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.List(ctx, limit, offset)
}

// UpdateTemplate updates a template
func (uc *CollageTemplateUseCase) UpdateTemplate(ctx context.Context, templateID uuid.UUID, name, filePath string) (*collage_template.CollageTemplate, error) {
	// テンプレートを取得
	template, err := uc.repo.FindByID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, collage_template.ErrTemplateNotFound
	}

	// 名前の更新
	if name != "" && name != template.Name() {
		if err := template.UpdateName(name); err != nil {
			return nil, err
		}
	}

	// ファイルパスの更新
	if filePath != "" && filePath != template.FilePath() {
		if err := template.UpdateFilePath(filePath); err != nil {
			return nil, err
		}
	}

	if err := uc.repo.Update(ctx, template); err != nil {
		return nil, err
	}

	return template, nil
}

// DeleteTemplate deletes a template
func (uc *CollageTemplateUseCase) DeleteTemplate(ctx context.Context, templateID uuid.UUID) error {
	// テンプレートが存在するかチェック
	template, err := uc.repo.FindByID(ctx, templateID)
	if err != nil {
		return err
	}
	if template == nil {
		return collage_template.ErrTemplateNotFound
	}

	return uc.repo.Delete(ctx, templateID)
}
