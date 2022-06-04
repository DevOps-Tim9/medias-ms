package repository

import (
	"medias-ms/src/entity"

	"gorm.io/gorm"
)

type IMediaRepository interface {
	Create(entity.Media) (entity.Media, error)
	Delete(uint)
	GetById(uint) (*entity.Media, error)
}

type MediaRepository struct {
	Database *gorm.DB
}

func (r MediaRepository) Create(media entity.Media) (entity.Media, error) {
	err := r.Database.Save(&media).Error

	return media, err
}

func (r MediaRepository) Delete(id uint) {
	r.Database.Unscoped().Delete(&entity.Media{}, id)
}

func (r MediaRepository) GetById(id uint) (*entity.Media, error) {
	var media = entity.Media{}

	err := r.Database.First(&media, id).Error

	return &media, err
}
