package server

import (
	"server/internal/executor"
	"server/internal/handler"
	"server/pkg/logger"
)

type Server struct {
	H *handler.Handler
}

func NewServer() *Server {

	execcpp := executor.NewExecutorCPP()
	lgr := logger.NewLogger()

	handler := handler.NewHandler(execcpp, lgr)

	return &Server{
		H: handler,
	}
}
