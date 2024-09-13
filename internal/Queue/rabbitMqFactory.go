package Queue

import (
	"github.com/streadway/amqp"
	"short-link/internal/Config"
	service_interface "short-link/internal/Core/Ports"
)

type Queue struct {
	Connection *amqp.Connection
	cfg        *Config.Config
	Service    service_interface.LinkServiceInterface
}

func CreateQueue(cfg *Config.Config) *Queue {

	queue := &Queue{
		Connection: CreateConnection(cfg),
		cfg:        cfg,
		Service:    nil,
	}

	return queue
}
