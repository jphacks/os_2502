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

func (uc *TemplatePartUseCase) CreateTemplatePart(
	ctx context.Context,
	templateID uuid.UUID,
	partNumber, positionX, positionY, width, height int,
	partName, description *string,
) (*template_part.TemplatePart, error) {
	// 同じテンプレートIDとパーツ番号の組み合わせが既に存在するかチェック
	existing, err := uc.repo.FindByTemplateIDAndPartNumber(ctx, templateID, partNumber)
	if err == nil && existing != nil {
		return nil, template_part.ErrDuplicatePartNumber
	}

	tp, err := template_part.NewTemplatePart(templateID, partNumber, positionX, positionY, width, height, partName, description)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, tp); err != nil {
		return nil, err
	}

	return tp, nil
}

func (uc *TemplatePartUseCase) GetTemplatePartByID(ctx context.Context, partID uuid.UUID) (*template_part.TemplatePart, error) {
	return uc.repo.FindByID(ctx, partID)
}

func (uc *TemplatePartUseCase) GetTemplatePartsByTemplateID(ctx context.Context, templateID uuid.UUID) ([]*template_part.TemplatePart, error) {
	return uc.repo.FindByTemplateID(ctx, templateID)
}

func (uc *TemplatePartUseCase) GetTemplatePartByTemplateIDAndPartNumber(ctx context.Context, templateID uuid.UUID, partNumber int) (*template_part.TemplatePart, error) {
	return uc.repo.FindByTemplateIDAndPartNumber(ctx, templateID, partNumber)
}

func (uc *TemplatePartUseCase) UpdateTemplatePartPosition(
	ctx context.Context,
	partID uuid.UUID,
	positionX, positionY, width, height int,
) (*template_part.TemplatePart, error) {
	tp, err := uc.repo.FindByID(ctx, partID)
	if err != nil {
		return nil, err
	}

	if err := tp.UpdatePosition(positionX, positionY, width, height); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, tp); err != nil {
		return nil, err
	}

	return tp, nil
}

func (uc *TemplatePartUseCase) UpdateTemplatePartName(
	ctx context.Context,
	partID uuid.UUID,
	partName *string,
) (*template_part.TemplatePart, error) {
	tp, err := uc.repo.FindByID(ctx, partID)
	if err != nil {
		return nil, err
	}

	tp.UpdatePartName(partName)

	if err := uc.repo.Update(ctx, tp); err != nil {
		return nil, err
	}

	return tp, nil
}

func (uc *TemplatePartUseCase) UpdateTemplatePartDescription(
	ctx context.Context,
	partID uuid.UUID,
	description *string,
) (*template_part.TemplatePart, error) {
	tp, err := uc.repo.FindByID(ctx, partID)
	if err != nil {
		return nil, err
	}

	tp.UpdateDescription(description)

	if err := uc.repo.Update(ctx, tp); err != nil {
		return nil, err
	}

	return tp, nil
}

func (uc *TemplatePartUseCase) DeleteTemplatePart(ctx context.Context, partID uuid.UUID) error {
	return uc.repo.Delete(ctx, partID)
}

func (uc *TemplatePartUseCase) DeleteTemplatePartsByTemplateID(ctx context.Context, templateID uuid.UUID) error {
	return uc.repo.DeleteByTemplateID(ctx, templateID)
}

func (uc *TemplatePartUseCase) ListTemplateParts(ctx context.Context, limit, offset int) ([]*template_part.TemplatePart, error) {
	return uc.repo.List(ctx, limit, offset)
}
