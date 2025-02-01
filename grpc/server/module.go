package server

import "go.uber.org/fx"

var Module = fx.Provide(
	NewStravaGrpcServer,
	NewGrpcServer,
)
