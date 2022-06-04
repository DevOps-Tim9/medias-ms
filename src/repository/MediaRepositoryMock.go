package repository

import (
	"errors"
	"medias-ms/src/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MediaRepositoryMock struct {
	mock.Mock
}

func (r MediaRepositoryMock) Create(media entity.Media) (entity.Media, error) {
	media.ID = 5
	return media, nil
}

func (r MediaRepositoryMock) Delete(id uint) {
}

func (r MediaRepositoryMock) GetById(id uint) (*entity.Media, error) {
	switch id {
	case 1:
		return &entity.Media{
			Model: gorm.Model{
				ID: 1,
			},
			Url: "/this/is/ful/media/path",
		}, nil
	case 2:
		return nil, errors.New("")
	}
	return nil, nil
}
