package main

import (
	"fmt"
	"medias-ms/src/config"
	config_db "medias-ms/src/config/db"
	"medias-ms/src/controller"
	"medias-ms/src/rabbitmq"
	"medias-ms/src/repository"
	"medias-ms/src/route"
	"medias-ms/src/service"
	"medias-ms/src/utils"
	"net/http"
	"os"

	"github.com/rs/cors"
	"gorm.io/gorm"
)

func main() {
	logger := utils.Logger()

	logger.Info("Connecting with DB")
	dataBase, _ := config_db.SetupDB()

	repositoryContainer := initializeRepositories(dataBase)
	serviceContainer := initializeServices(repositoryContainer)
	controllerContainer := initializeControllers(serviceContainer)

	router := route.SetupRoutes(controllerContainer)

	fileServer := http.FileServer(http.Dir("./static/"))

	router.PathPrefix("/static/").Handler(http.StripPrefix("/api/static/", fileServer))

	port := os.Getenv("SERVER_PORT")

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	logger.Info("Connecting on RabbitMq")

	rabbit := rabbitmq.RMQConsumer{
		ConnectionString: amqpServerURL,
		MediaService:     serviceContainer.MediaService,
	}

	channel, _ := rabbit.StartRabbitMQ()

	defer channel.Close()

	messages, _ := channel.Consume(
		"DeleteImageOnMedias-MS",          // queue name
		"DeleteImageOnMedias-MS-consumer", // consumer
		true,                              // auto-ack
		false,                             // exclusive
		false,                             // no local
		false,                             // no wait
		nil,                               // arguments
	)

	go rabbit.Worker(messages)

	logger.Info("Starting server")

	http.ListenAndServe(fmt.Sprintf(":%s", port), cors.AllowAll().Handler(router))
}

func initializeControllers(serviceContainer config.ServiceContainer) config.ControllerContainer {
	mediaController := controller.NewMediaController(serviceContainer.MediaService)

	container := config.NewControllerContainer(
		mediaController,
	)

	return container
}

func initializeServices(repositoryContainer config.RepositoryContainer) config.ServiceContainer {
	mediaService := service.MediaService{
		MediaRepository: repositoryContainer.MediaRepository,
		Logger:          utils.Logger(),
	}

	container := config.NewServiceContainer(
		mediaService,
	)

	return container
}

func initializeRepositories(dataBase *gorm.DB) config.RepositoryContainer {
	mediaRepository := repository.MediaRepository{Database: dataBase}

	container := config.NewRepositoryContainer(
		mediaRepository,
	)

	return container
}
