package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/template_part"
)

type TemplatePartUseCase struct {
	repo template_part.Repository
}

func NewTemplatePartUseCase(repo template_part.Repository) *TemplatePartUseCase {
	return &TemplatePartUseCase{repo: repo}
}

func (uc *TemplatePartUseCase) CreateTemplatePart(ctx context.Context, templateID uuid.UUID, partNumber, positionX, positionY, width, height int) (*template_part.TemplatePart, error) {
	// テンプレートパーツを作成
	part, err := template_part.NewTemplatePart(templateID, partNumber, positionX, positionY, width, height)
	if err != nil {
		return nil, err
	}

	// リポジトリに保存
	if err := uc.repo.Save(ctx, part); err != nil {
		return nil, err
	}

	return part, nil
}

func (uc *TemplatePartUseCase) GetTemplatePartByID(ctx context.Context, partID uuid.UUID) (*template_part.TemplatePart, error) {
	return uc.repo.FindByID(ctx, partID)
}

func (uc *TemplatePartUseCase) GetTemplatePartsByTemplateID(ctx context.Context, templateID uuid.UUID) ([]*template_part.TemplatePart, error) {
	return uc.repo.FindByTemplateID(ctx, templateID)
}

func (uc *TemplatePartUseCase) UpdatePartName(ctx context.Context, partID uuid.UUID, name string) (*template_part.TemplatePart, error) {
	part, err := uc.repo.FindByID(ctx, partID)
	if err != nil {
		return nil, err
	}

	part.SetPartName(name)

	if err := uc.repo.Save(ctx, part); err != nil {
		return nil, err
	}

	return part, nil
}

func (uc *TemplatePartUseCase) UpdatePartDescription(ctx context.Context, partID uuid.UUID, description string) (*template_part.TemplatePart, error) {
	part, err := uc.repo.FindByID(ctx, partID)
	if err != nil {
		return nil, err
	}

	part.SetDescription(description)

	if err := uc.repo.Save(ctx, part); err != nil {
		return nil, err
	}

	return part, nil
}

func (uc *TemplatePartUseCase) UpdatePartPosition(ctx context.Context, partID uuid.UUID, x, y int) (*template_part.TemplatePart, error) {
	part, err := uc.repo.FindByID(ctx, partID)
	if err != nil {
		return nil, err
	}

	if err := part.UpdatePosition(x, y); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, part); err != nil {
		return nil, err
	}

	return part, nil
}

func (uc *TemplatePartUseCase) UpdatePartDimensions(ctx context.Context, partID uuid.UUID, width, height int) (*template_part.TemplatePart, error) {
	part, err := uc.repo.FindByID(ctx, partID)
	if err != nil {
		return nil, err
	}

	if err := part.UpdateDimensions(width, height); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, part); err != nil {
		return nil, err
	}

	return part, nil
}

func (uc *TemplatePartUseCase) DeleteTemplatePart(ctx context.Context, partID uuid.UUID) error {
	return uc.repo.Delete(ctx, partID)
}
