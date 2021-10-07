package interceptors

import (
	"context"
	"time"

	"github.com/quyenphamkhac/gmd-productsrv/config"
	"github.com/quyenphamkhac/gmd-productsrv/internal/logger"
	"github.com/quyenphamkhac/gmd-productsrv/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type InterceptorManager struct {
	logger logger.Logger
	cfg    *config.Config
	metr   metrics.Metrics
}

func NewInterceptorManager(logger logger.Logger, cfg *config.Config, metr metrics.Metrics) *InterceptorManager {
	return &InterceptorManager{
		logger: logger,
		cfg:    cfg,
		metr:   metr,
	}
}

func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	resp, err = handler(ctx, req)
	im.logger.Infof("Method: %s, Time: %v, Metadata: %v, Err: %v", info.FullMethod, time.Since(start), md, err)
	return resp, err
}
