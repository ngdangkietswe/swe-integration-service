package kafka

import (
	"github.com/ngdangkietswe/swe-integration-service/kafka/consumer"
	"github.com/ngdangkietswe/swe-integration-service/kafka/handler"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	// Handlers
	handler.NewCdcAuthUsersHandler,

	// Consumers
	consumer.NewCdcAuthUsersConsumer,
)
