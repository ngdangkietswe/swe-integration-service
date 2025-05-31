package logger

import (
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"github.com/ngdangkietswe/swe-go-common-shared/logger"
	"go.uber.org/fx"
)

func NewZapLogger() (*logger.Logger, error) {
	instance, err := logger.NewLogger(
		config.GetString("APP_NAME", "swe-integration"),
		config.GetString("APP_ENV", "dev"),
		"debug",
		"",
	)

	if err != nil {
		return nil, err
	}

	return instance, nil
}

var Module = fx.Provide(NewZapLogger)
