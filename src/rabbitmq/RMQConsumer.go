package rabbitmq

import (
	"context"
	"encoding/json"
	"medias-ms/src/dto"
	"medias-ms/src/service"

	"github.com/streadway/amqp"
)

type RMQConsumer struct {
	ConnectionString string
	MediaService     service.IMediaService
}

func (r RMQConsumer) StartRabbitMQ() (*amqp.Channel, error) {
	connectRabbitMQ, _ := amqp.Dial(r.ConnectionString)

	channelRabbitMQ, _ := connectRabbitMQ.Channel()

	err := channelRabbitMQ.ExchangeDeclare(
		"DeleteImageOnMedias-MS-exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	queue, err := channelRabbitMQ.QueueDeclare(
		"DeleteImageOnMedias-MS",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	err = channelRabbitMQ.QueueBind(
		queue.Name,
		"DeleteImageOnMedias-MS-routing-key",
		"DeleteImageOnMedias-MS-exchange",
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	err = channelRabbitMQ.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		return nil, err
	}

	return channelRabbitMQ, nil
}

func (r RMQConsumer) HandleDeleteImage(message []byte) {
	var mediaDto dto.MediaDto

	json.Unmarshal([]byte(message), &mediaDto)

	r.MediaService.Delete(mediaDto.Id, context.TODO())
}

func (r RMQConsumer) Worker(messages <-chan amqp.Delivery) {
	for delivery := range messages {
		r.HandleDeleteImage(delivery.Body)
	}
}
