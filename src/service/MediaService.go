package service

import (
	"context"
	"io"
	"medias-ms/src/entity"
	"medias-ms/src/repository"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type IMediaService interface {
	Save(*multipart.FileHeader, context.Context) (entity.Media, error)
	Delete(uint, context.Context) error
	GetById(uint, context.Context) (*entity.Media, error)
}

type MediaService struct {
	MediaRepository repository.IMediaRepository
	Logger          *logrus.Entry
}

func (s MediaService) GetById(id uint, ctx context.Context) (*entity.Media, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Get media by id")

	defer span.Finish()

	media, err := s.MediaRepository.GetById(id, ctx)

	if err != nil {
		return nil, err
	}

	return media, nil
}

func (s MediaService) Save(file *multipart.FileHeader, ctx context.Context) (entity.Media, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Create media")

	defer span.Finish()

	s.Logger.Info("Saving media in file system")

	url := s.SaveFile(file)

	media := entity.Media{
		Url: url,
	}

	s.Logger.Info("Saving media in database.")

	media, error := s.MediaRepository.Create(media, ctx)

	return media, error
}

func (s MediaService) Delete(id uint, ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Delete media by id")

	defer span.Finish()

	s.Logger.Info("Finding media by id")

	media, err := s.GetById(id, ctx)

	if err != nil {
		s.Logger.Error("Media doesn't exist")

		return err
	}

	s.Logger.Info("Started deleting media from file system")

	os.Remove("./static/images/" + strings.Split(media.Url, "/")[5])

	s.Logger.Info("Started deleting media from DB")

	s.MediaRepository.Delete(id, ctx)

	return nil
}

func (s MediaService) SaveFile(fileHeader *multipart.FileHeader) string {
	s.Logger.Info("Saving media in file system")

	fileName, _ := uuid.NewV4()

	destinationFilePath, uriPathToImage := s.createDestinationFilePathAndUriPath(fileName.String())

	file, _ := fileHeader.Open()
	defer file.Close()

	dir := path.Dir(destinationFilePath)
	os.MkdirAll(dir, os.ModeDir|0700)

	var output *os.File
	output, _ = os.OpenFile(destinationFilePath, os.O_CREATE|os.O_WRONLY, 0600)

	io.Copy(output, file)

	output.Close()

	s.Logger.Info("Media saved in file system")

	return "/medias-ms/" + uriPathToImage
}

func (s MediaService) createDestinationFilePathAndUriPath(fileName string) (string, string) {
	fullFileName := fileName + ".png"

	destinationFilePath := "./static/images/" + fullFileName

	uriPathToImage := "api/static/images/" + fullFileName

	return destinationFilePath, uriPathToImage
}
