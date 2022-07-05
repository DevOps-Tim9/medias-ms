package service

import (
	"context"
	"fmt"
	"medias-ms/src/entity"
	"medias-ms/src/repository"
	"medias-ms/src/utils"
	"mime/multipart"
	"net/textproto"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MediaServiceIntegrationTestSuite struct {
	suite.Suite
	service MediaService
	db      *gorm.DB
	medias  []entity.Media
}

func (suite *MediaServiceIntegrationTestSuite) SetupSuite() {
	host := os.Getenv("DATABASE_DOMAIN")
	user := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	name := os.Getenv("DATABASE_SCHEMA")
	port := os.Getenv("DATABASE_PORT")

	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		name,
		port,
	)

	db, _ := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	db.AutoMigrate(&entity.Media{Tbl: "medias"})

	mediaRepository := repository.MediaRepository{Database: db}

	suite.db = db
	suite.service = MediaService{
		MediaRepository: mediaRepository,
		Logger:          utils.Logger(),
	}

	suite.medias = []entity.Media{
		{
			Model: gorm.Model{
				ID:        100,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Url: "http://localhost:8082/static/media/some_name1.png",
		},
		{
			Model: gorm.Model{
				ID:        200,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Url: "http://localhost:8082/static/media/some_name2.png",
		},
	}

	tx := suite.db.Begin()

	tx.Create(&suite.medias[0])
	tx.Create(&suite.medias[1])

	tx.Commit()
}

func TestMediaServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(MediaServiceIntegrationTestSuite))
}

func (suite *MediaServiceIntegrationTestSuite) TestIntegrationMediaService_GetById_MediaDoesNotExist() {
	media, err := suite.service.GetById(5, context.TODO())

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), media)
}

func (suite *MediaServiceIntegrationTestSuite) TestIntegrationMediaService_GetById_MediaDoesExist() {
	id := uint(100)

	media, err := suite.service.GetById(id, context.TODO())

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), media)
	assert.Equal(suite.T(), media.ID, id)
}

func (suite *MediaServiceIntegrationTestSuite) TestIntegrationMediaService_Delete_MediaDoesExist() {
	id := uint(200)

	err := suite.service.Delete(id, context.TODO())

	assert.Nil(suite.T(), err)
}

func (suite *MediaServiceIntegrationTestSuite) TestIntegrationMediaService_Delete_MediaDoesNotExist() {
	id := uint(500)

	err := suite.service.Delete(id, context.TODO())

	assert.NotNil(suite.T(), err)
}

func (suite *MediaServiceIntegrationTestSuite) TestIntegrationMediaService_Create_Success() {
	id := uint(1)

	file := &multipart.FileHeader{
		Filename: "a",
		Header:   textproto.MIMEHeader{},
		Size:     0,
	}

	media, err := suite.service.Save(file, context.TODO())

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), media.ID, id)
}
