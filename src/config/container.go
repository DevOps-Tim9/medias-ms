package config

import (
	"medias-ms/src/controller"
	"medias-ms/src/repository"
	"medias-ms/src/service"
)

type ControllerContainer struct {
	MediaController controller.MediaController
}

type ServiceContainer struct {
	MediaService service.IMediaService
}

type RepositoryContainer struct {
	MediaRepository repository.IMediaRepository
}

func NewControllerContainer(
	mediaController controller.MediaController,

) ControllerContainer {
	return ControllerContainer{
		MediaController: mediaController,
	}
}

func NewServiceContainer(mediaService service.IMediaService) ServiceContainer {
	return ServiceContainer{
		MediaService: mediaService,
	}
}

func NewRepositoryContainer(
	mediaRepository repository.IMediaRepository,

) RepositoryContainer {
	return RepositoryContainer{
		MediaRepository: mediaRepository,
	}
}
