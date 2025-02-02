package grpc

import (
	"github.com/ngdangkietswe/swe-integration-service/grpc/server"
	"github.com/ngdangkietswe/swe-integration-service/grpc/service"
	"github.com/ngdangkietswe/swe-integration-service/grpc/validator"
	"go.uber.org/fx"
)

var Module = fx.Options(
	validator.Module,
	service.Module,
	server.Module,
)
