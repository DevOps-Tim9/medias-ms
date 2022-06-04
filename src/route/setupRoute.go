package route

import (
	"medias-ms/src/config"

	"github.com/gorilla/mux"
)

func SetupRoutes(container config.ControllerContainer) *mux.Router {
	route := mux.NewRouter()

	routerWithApiAsPrefix := route.PathPrefix("/api").Subrouter()

	routerWithApiAsPrefix.HandleFunc("/medias", container.MediaController.Upload).Methods("POST")
	routerWithApiAsPrefix.HandleFunc("/medias/{id}", container.MediaController.Delete).Methods("DELETE")
	routerWithApiAsPrefix.HandleFunc("/medias/{id}", container.MediaController.GetById).Methods("GET")

	return routerWithApiAsPrefix
}
