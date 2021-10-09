package server

import (
	"github.com/quyenphamkhac/gmd-productsrv/config"
	"github.com/quyenphamkhac/gmd-productsrv/internal/logger"
	"github.com/streadway/amqp"
)

type Server struct {
	logger   logger.Logger
	cfg      *config.Config
	amqpConn *amqp.Connection
}

func NewProductServer(logger logger.Logger, cfg *config.Config, amqpConn *amqp.Connection) *Server {
	return &Server{
		amqpConn: amqpConn,
		logger:   logger,
		cfg:      cfg,
	}
}

func (s *Server) Run() error {
	return nil
}
