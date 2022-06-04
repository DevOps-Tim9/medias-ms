package controller

import (
	"encoding/json"
	"medias-ms/src/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v8"
)

type MediaController struct {
	MediaService service.IMediaService
	validate     *validator.Validate
}

func NewMediaController(mediaService service.IMediaService) MediaController {
	config := &validator.Config{TagName: "validate"}

	return MediaController{MediaService: mediaService, validate: validator.New(config)}
}

func (c MediaController) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	files := r.MultipartForm.File["files"]

	media, error := c.MediaService.Save(files[0])

	if error != nil {
		handleMediaError(error, w)

		return
	}

	payload, _ := json.Marshal(media)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(payload))
}

func (c MediaController) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	c.MediaService.Delete(uint(id))

	w.WriteHeader(http.StatusNoContent)
}

func (c MediaController) GetById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	media, err := c.MediaService.GetById(uint(id))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	payload, _ := json.Marshal(media)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func handleMediaError(error error, w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusInternalServerError)

	return w
}
