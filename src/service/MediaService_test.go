package service

import (
	"errors"
	"medias-ms/src/entity"
	"medias-ms/src/repository"
	"mime/multipart"
	"net/textproto"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type MediaServiceUnitTestSuite struct {
	suite.Suite
	mediaRepositoryMock *repository.MediaRepositoryMock
	service             MediaService
}

func TestMediaServiceUnitTestSuite(t *testing.T) {
	suite.Run(t, new(MediaServiceUnitTestSuite))
}

func (suite *MediaServiceUnitTestSuite) SetupSuite() {
	suite.mediaRepositoryMock = new(repository.MediaRepositoryMock)

	suite.service = MediaService{MediaRepository: suite.mediaRepositoryMock}
}

func (suite *MediaServiceUnitTestSuite) TestNewPostService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *MediaServiceUnitTestSuite) TestMediaService_GetById_MediaNotExist() {
	err := errors.New("")

	media, returnedErr := suite.service.GetById(2)

	assert.Nil(suite.T(), media, "Media not nil")
	assert.NotNil(suite.T(), returnedErr, "Error is nil")
	assert.Equal(suite.T(), err, returnedErr, "Not equal")
}

func (suite *MediaServiceUnitTestSuite) TestMediaService_GetById_MediaExist() {
	media := entity.Media{
		Model: gorm.Model{
			ID: 1,
		},
		Url: "/this/is/ful/media/path",
	}

	returnedMedia, err := suite.service.GetById(1)

	assert.Nil(suite.T(), err, "Error not nil")
	assert.NotNil(suite.T(), returnedMedia, "Media is nil")
	assert.Equal(suite.T(), media.ID, returnedMedia.ID, "Not equal")
	assert.Equal(suite.T(), media.Url, returnedMedia.Url, "Not equal")
}

func (suite *MediaServiceUnitTestSuite) TestMediaService_Delete_MediaNotExist() {
	err := suite.service.Delete(2)

	assert.NotNil(suite.T(), err, "Error is nil")
}

func (suite *MediaServiceUnitTestSuite) TestMediaService_Delete_MediaExist() {
	err := suite.service.Delete(1)

	assert.Nil(suite.T(), err, "Error is not nil")
}

func (suite *MediaServiceUnitTestSuite) TestMediaService_Create_ReturnsMedia() {
	id := uint(5)

	file := &multipart.FileHeader{
		Filename: "a",
		Header:   textproto.MIMEHeader{},
		Size:     0,
	}

	media, err := suite.service.Save(file)

	assert.Nil(suite.T(), err, "Error is not nil")
	assert.NotNil(suite.T(), media, "Media is nil")
	assert.Equal(suite.T(), id, media.ID, "Ids not equal")

	os.RemoveAll("static")
}
