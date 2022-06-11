package controller

import (
	"encoding/json"
	"medias-ms/src/service"
	"medias-ms/src/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v8"
)

type MediaController struct {
	MediaService service.IMediaService
	validate     *validator.Validate
	logger       *logrus.Entry
}

func NewMediaController(mediaService service.IMediaService) MediaController {
	config := &validator.Config{TagName: "validate"}
	logger := utils.Logger()

	return MediaController{MediaService: mediaService, validate: validator.New(config), logger: logger}
}

func (c MediaController) Upload(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Uploading media request received")

	r.ParseMultipartForm(32 << 20)

	files := r.MultipartForm.File["files"]

	media, error := c.MediaService.Save(files[0])

	if error != nil {
		handleMediaError(error, w)

		c.logger.Error("Error occured in uploading media")

		return
	}

	payload, _ := json.Marshal(media)

	c.logger.Info("Media uploaded successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))
}

func (c MediaController) Delete(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Deleting media request received")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		c.logger.Error("Error occured in deleting media")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	c.MediaService.Delete(uint(id))

	c.logger.Info("Media deleted successfully.")

	w.WriteHeader(http.StatusNoContent)
}

func (c MediaController) GetById(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Finding media by id request received")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		c.logger.Error("Error occured in finding media")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	media, err := c.MediaService.GetById(uint(id))

	if err != nil {
		c.logger.Error("Media not found")

		w.WriteHeader(http.StatusNotFound)

		return
	}

	payload, _ := json.Marshal(media)

	c.logger.Info("Media found successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func handleMediaError(error error, w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusInternalServerError)

	return w
}
