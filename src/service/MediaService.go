package service

import (
	"io"
	"medias-ms/src/entity"
	"medias-ms/src/repository"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/gofrs/uuid"
)

type IMediaService interface {
	Save(*multipart.FileHeader) (entity.Media, error)
	Delete(uint) error
	GetById(uint) (*entity.Media, error)
}

type MediaService struct {
	MediaRepository repository.IMediaRepository
}

func (s MediaService) GetById(id uint) (*entity.Media, error) {
	media, err := s.MediaRepository.GetById(id)

	if err != nil {
		return nil, err
	}

	return media, nil
}

func (s MediaService) Save(file *multipart.FileHeader) (entity.Media, error) {
	url := s.SaveFile(file)

	media := entity.Media{
		Url: url,
	}

	media, error := s.MediaRepository.Create(media)

	return media, error
}

func (s MediaService) Delete(id uint) error {
	media, err := s.GetById(id)

	if err != nil {
		return err
	}

	os.Remove("./static/images/" + strings.Split(media.Url, "/")[5])

	s.MediaRepository.Delete(id)

	return nil
}

func (s MediaService) SaveFile(fileHeader *multipart.FileHeader) string {
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

	return "/medias-ms/" + uriPathToImage
}

func (s MediaService) createDestinationFilePathAndUriPath(fileName string) (string, string) {
	fullFileName := fileName + ".png"

	destinationFilePath := "./static/images/" + fullFileName

	uriPathToImage := "api/static/images/" + fullFileName

	return destinationFilePath, uriPathToImage
}
