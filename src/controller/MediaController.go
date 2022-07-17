package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"medias-ms/src/dto"
	"medias-ms/src/service"
	"medias-ms/src/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Handle /api/medias")

	defer span.Finish()

	c.logger.Info("Uploading media request received")

	r.ParseMultipartForm(32 << 20)

	files := r.MultipartForm.File["files"]

	media, error := c.MediaService.Save(files[0], ctx)

	if error != nil {
		handleMediaError(error, w)

		c.logger.Error("Error occured in uploading media")

		AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), "New media uploaded unsuccessfully")

		return
	}

	payload, _ := json.Marshal(media)

	c.logger.Info("Media uploaded successfully")

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), "New media uploaded successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))
}

func (c MediaController) Delete(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Handle /api/medias/{id}")

	defer span.Finish()

	c.logger.Info("Deleting media request received")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		c.logger.Error("Error occured in deleting media")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	c.MediaService.Delete(uint(id), ctx)

	c.logger.Info("Media deleted successfully.")

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("Media with id %d deleted successfully", id))

	w.WriteHeader(http.StatusNoContent)
}

func (c MediaController) GetById(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Handle /api/medias/{id}")

	defer span.Finish()

	c.logger.Info("Finding media by id request received")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		c.logger.Error("Error occured in finding media")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	media, err := c.MediaService.GetById(uint(id), ctx)

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

func AddSystemEvent(time string, message string) error {
	logger := utils.Logger()
	event := dto.EventRequestDTO{
		Timestamp: time,
		Message:   message,
	}

	b, _ := json.Marshal(&event)
	endpoint := os.Getenv("EVENTS_MS")
	logger.Info("Sending system event to events-ms")
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")

	_, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Debug("Error happened during sending system event")
		return err
	}

	return nil
}
