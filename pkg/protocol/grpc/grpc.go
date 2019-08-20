package grpc

import (
	grpc_middleware_logging_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
)

// SetGRPCLogger replace logger that used by grpc
func SetGRPCLogger(logger *zap.Logger) {
	grpc_middleware_logging_zap.ReplaceGrpcLogger(logger)
}
