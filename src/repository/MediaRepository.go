package repository

import (
	"context"
	"medias-ms/src/entity"

	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

type IMediaRepository interface {
	Create(entity.Media, context.Context) (entity.Media, error)
	Delete(uint, context.Context)
	GetById(uint, context.Context) (*entity.Media, error)
}

type MediaRepository struct {
	Database *gorm.DB
}

func (r MediaRepository) Create(media entity.Media, ctx context.Context) (entity.Media, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Create media")

	defer span.Finish()

	err := r.Database.Save(&media).Error

	return media, err
}

func (r MediaRepository) Delete(id uint, ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Delete media by id")

	defer span.Finish()

	r.Database.Unscoped().Delete(&entity.Media{}, id)
}

func (r MediaRepository) GetById(id uint, ctx context.Context) (*entity.Media, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Get media by id")

	defer span.Finish()

	var media = entity.Media{}

	err := r.Database.First(&media, id).Error

	return &media, err
}
